package loadbalancer

import (
	"flux/utils"
	"flux/datastore/redis"
	"net/http"
	"github.com/joho/godotenv"
	"os"

)

var configData *utils.Config

func Initialize(config *utils.Config) {
	godotenv.Load()
	configData = config
	redis := redis.Redis{}
	redis.Connect(config.Redis_host,config.Load_redis_dump,config.Redis_dump_interval)
	redis.InitServerList(&config.Initial_servers)

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request){
		if req.URL.Path=="/add-server" || req.URL.Path=="/remove-server"{
			return 
		}


	})
	
	http.HandleFunc("/add-server",func(res http.ResponseWriter, req *http.Request){
		queryParams := req.URL.Query()
		server := queryParams.Get("server")
		code := queryParams.Get("code")
		if code != os.Getenv("SECRET"){
			return
		}

		redis.AddMember(server)


	})


	http.HandleFunc("/remove-server", func(res http.ResponseWriter, req *http.Request){
		queryParams := req.URL.Query()
		server := queryParams.Get("server")
		code := queryParams.Get("code")
		if code != os.Getenv("SECRET"){
			return
		}

		redis.RemoveMember(server)
	})
	



}

func Start() {
	http.ListenAndServe(configData.App_host,nil)
	
}