package utils

import (
	"encoding/json"
	"fmt"
	"os"

)

type Config struct {
	App_host string
	Redis_host string
	Algorithm string
	Timeout timeoutConfig
	Server_retries int
	Load_redis_dump bool
	Health_check_interval int
	Redis_dump_interval string
	Initial_servers []string
}

type timeoutConfig struct {
	Health_check int 
	Server_response int 

}

func ParseJSON() *Config{

	configData, err := os.ReadFile("./config/config.json")
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return nil
	}

	// Parse JSON into a map
	var config map[string]interface{}
	err = json.Unmarshal(configData, &config)
	if err != nil {
		fmt.Println("Error parsing config:", err)
		return nil
	}

	app_host := config["appHost"].(string)
	redis_host := config["redisHost"].(string)
	algorithm := config["algorithm"].(string)
	server_retries := int(config["serverRetries"].(float64))
	load_redis_dump := config["loadRedisDump"].(bool)
	redis_dump_interval := config["redisDumpInterval"].(string)
	health_check_interval := int(config["healthCheckInterval"].(float64))

	interface_slice := config["initialServers"].([]interface{})
	initial_servers := make([]string,len(interface_slice))
	for i, element := range interface_slice{
		initial_servers[i]=element.(string)
	}

	timeout := config["timeout"].(map[string]interface{})
	health_check_timeout  := int(timeout["healthCheck"].(float64))
	server_response_timeout  := int(timeout["serverResponse"].(float64))
	
	configStruct := Config{
	App_host: app_host,
	Redis_host: redis_host,
	Algorithm: algorithm,
	Server_retries: server_retries,
	Load_redis_dump: load_redis_dump,
	Redis_dump_interval: redis_dump_interval,
	Health_check_interval: health_check_interval,
	Timeout: timeoutConfig{
		Health_check: health_check_timeout, 
		Server_response: server_response_timeout,
	},
	Initial_servers: initial_servers,
	}
	
	fmt.Println(configStruct)
	return &configStruct

 
}

