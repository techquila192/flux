package loadbalancer

import (
	"flux/utils"
	"flux/datastore/redis"
	"net/http"

)

var configData *utils.Config

func Initialize(config *utils.Config) {
	configData = config
	redis := redis.Redis{}
	redis.Connect(config.Redis_host,config.Load_redis_dump,config.Redis_dump_interval)
	redis.InitServerList(&config.Initial_servers)
	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request){

	})
	
	//routes
	//middleware


}

func Start() {
	http.ListenAndServe(configData.App_host,nil)
	
}