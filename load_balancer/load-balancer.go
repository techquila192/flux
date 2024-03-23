package loadbalancer

import (
	"flux/utils"
	"flux/algorithm"
	"flux/datastore/redis"
	"net/http"
	"github.com/joho/godotenv"
	"sync"
	"time"
	"fmt"
	"os"

)

var configData *utils.Config
var redisInstance redis.Redis
var visited map[string] bool
var visited_mu sync.Mutex 
var serviceClient http.Client


func Start() {
	//start health check in other goroutine
	ticker := time.NewTicker(time.Duration(configData.Health_check_interval) * time.Second)
	go func(){
		<- ticker.C
		utils.HealthCheck(redisInstance.GetServers(),&redisInstance.Mu,configData.Timeout.Health_check)
	}()
	
	http.ListenAndServe(configData.App_host,nil)
	
}


func Initialize(config *utils.Config) {
	godotenv.Load()
	serviceClient = http.Client{
        Timeout: time.Duration(config.Timeout.Server_response) * time.Second, 
    }
	configData = config
	redisInstance = redis.Redis{}
	visited_mu = sync.Mutex{}
	visited = make(map[string]bool)
	redisInstance.Connect(config.Redis_host,config.Load_redis_dump,config.Redis_dump_interval)
	redisInstance.InitServerList(&config.Initial_servers)

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request){
		if req.URL.Path=="/add-server" || req.URL.Path=="/remove-server"{
			return 
		}

		sent := false

		for sent!=true	{
		server := getServer()  
		newRequest, err := http.NewRequest(req.Method, req.URL.Scheme + "://" + server.GetHost() + req.URL.Path, req.Body)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}
		server.IncrementConnections()
		newResponse, err := serviceClient.Do(newRequest)
		if err != nil {
			//check for timeout then retry else look for new 
		} else {
			sent = true
		}
		server.DecrementConnections()
		}



	})
	
	http.HandleFunc("/add-server",func(res http.ResponseWriter, req *http.Request){
		queryParams := req.URL.Query()
		server := queryParams.Get("server")
		code := queryParams.Get("code")
		if code != os.Getenv("SECRET"){
			return
		}

		redisInstance.AddMember(server)


	})


	http.HandleFunc("/remove-server", func(res http.ResponseWriter, req *http.Request){
		queryParams := req.URL.Query()
		server := queryParams.Get("server")
		code := queryParams.Get("code")
		if code != os.Getenv("SECRET"){
			return
		}

		redisInstance.RemoveMember(server)
	})
	
}



func getServer() *redis.Server{

	if configData.Algorithm == "round-robin"{
		return algorithm.Round_robin(redisInstance.GetServers(),&redisInstance.Mu,&visited,&visited_mu)
	}
}