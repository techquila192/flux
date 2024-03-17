package utils

import (
	"flux/datastore/redis"
	"sync"
	"net/http"
	"fmt"
	"time"
)

func HealthCheck(servers *map[string]*redis.Server, mu *sync.RWMutex, timeout int) {
	mu.Lock()
	defer mu.Unlock()
	var wg sync.WaitGroup
	wg.Add(len(*servers))
	client := &http.Client{
        Timeout: time.Duration(timeout) * time.Second, // Set timeout in seconds
    }
	for _, server := range *servers {
		go pingNode(client,server,&wg)
	}

	wg.Wait()

}

func pingNode(client *http.Client, server *redis.Server, wg *sync.WaitGroup) {
	
	defer wg.Done()
	response, err := client.Get("http://"+server.GetHost())
    if err != nil {
        fmt.Println("Error:", err)
		server.SetAliveState(false)
        return
    }
	defer response.Body.Close()
	
	if response.StatusCode == http.StatusOK {
        server.SetAliveState(true)
		fmt.Println(server.GetHost(),"is healthy")
    } else {
		fmt.Println(server.GetHost(),"is dead")
        server.SetAliveState(false)
    }
}