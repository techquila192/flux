package utils

import (
	"flux/datastore/redis"
	"sync"
	"net/http"
	"fmt"
)

func HealthCheck(servers *map[string]*redis.Server, mu *sync.RWMutex, timeout int) {
	mu.Lock()
	defer mu.Unlock()
	var wg sync.WaitGroup
	wg.Add(len(*servers))
	for _, server := range *servers {
		go pingNode(server,&wg)
	}

	wg.Wait()

}

func pingNode(server *redis.Server, wg *sync.WaitGroup) {
	defer wg.Done()
	response, err := http.Get(server.GetHost())
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer response.Body.Close()
	if response.StatusCode == http.StatusOK {
        server.SetAliveState(true)
    } else {
        server.SetAliveState(false)
    }
}