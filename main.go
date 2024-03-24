package main

import (
	"flux/utils"
	"flux/load_balancer"
)



func main(){

	config := utils.ParseJSON()
	loadbalancer.Initialize(config)
	loadbalancer.Start()
	
	
	
}
