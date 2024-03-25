package main

import (
	"flux/utils"
	"flux/load_balancer"
	"os"
	"os/signal"
    "syscall"
)



func main(){

	config := utils.ParseJSON()
	loadbalancer.Initialize(config)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
         <-sigCh
        loadbalancer.Close()
		os.Exit(0)    
    }()
	loadbalancer.Start()
	
	
}
