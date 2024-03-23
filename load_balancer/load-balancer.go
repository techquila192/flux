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
	"io"
	"os"

)

var configData *utils.Config
var redisInstance redis.Redis
var visited map[string] bool
var visited_mu sync.Mutex 
var serviceClient http.Client


func Start() {
	//start health check 
	ticker := time.NewTicker(time.Duration(configData.Health_check_interval) * time.Second)
	go func(){
		<- ticker.C
		utils.HealthCheck(redisInstance.GetServers(),&redisInstance.Mu,configData.Timeout.Health_check)
	}()
	
	//server startup
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

	http.HandleFunc("/", balancer)	
	http.HandleFunc("/add-server",addServer)
	http.HandleFunc("/remove-server", removeServer)
	
}

func balancer(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path=="/add-server" || req.URL.Path=="/remove-server"{
		return 
	}

	sent := false
	var newResponse *http.Response 

	for !sent	{
		server := getServer()  
		newRequest, err := http.NewRequest(req.Method, req.URL.Scheme + "://" + server.GetHost() + req.URL.Path, req.Body)
		copyRequestHeaders(req,newRequest)
		if err != nil {
			fmt.Println("Error creating request:", err)
			return
		}

		server.IncrementConnections()
		newResponse, err = serviceClient.Do(newRequest)
		server.DecrementConnections()
		if err != nil {
			//check for timeout then retry else look for new 
			if os.IsTimeout(err){
				for i:=1;i<=configData.Server_retries;i++{
					server.IncrementConnections()
					retryResponse, er := serviceClient.Do(newRequest)
					server.DecrementConnections()
					if er != nil{
						newResponse = retryResponse
						sent = true
						break
					}
				}
			} 
		} else {
			sent = true
		}
	
	}
	copyResponse(newResponse, res)

}

func copyRequestHeaders(src *http.Request, dest *http.Request) {
	for key, values := range src.Header {
		for _, value := range values {
			dest.Header.Add(key, value)
		}
	}
}

func copyResponse(resp *http.Response, w http.ResponseWriter) {
    // Copy headers from the response to the writer
    for key, values := range resp.Header {
        for _, value := range values {
            w.Header().Add(key, value)
        }
    }

    // Set status code
    w.WriteHeader(resp.StatusCode)

    // Copy response body to the writer
    _, err := io.Copy(w, resp.Body)
    if err != nil {
        fmt.Println("Error copying response body:", err)
    }
}

func addServer(res http.ResponseWriter, req *http.Request){
	queryParams := req.URL.Query()
	server := queryParams.Get("server")
	code := queryParams.Get("code")
	if code != os.Getenv("SECRET"){
		return
	}

	redisInstance.AddMember(server)


}

func removeServer(res http.ResponseWriter, req *http.Request){
	queryParams := req.URL.Query()
	server := queryParams.Get("server")
	code := queryParams.Get("code")
	if code != os.Getenv("SECRET"){
		return
	}

	redisInstance.RemoveMember(server)
}


func getServer() *redis.Server{

	if configData.Algorithm == "round-robin"{
		return algorithm.Round_robin(redisInstance.GetServers(),&redisInstance.Mu,&visited,&visited_mu)
	}

	return &redis.Server{}
}