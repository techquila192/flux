package algorithm

import (
	"flux/datastore/redis"
	"sync"
	"fmt"
)


type ServerHeap struct{
	items []string
	scheme string 
	mu sync.Mutex 
	serverMap *map[string] interface{}
}

func (sh *ServerHeap) Len() int { return len(sh.items) }

func (sh *ServerHeap) Less(i, j int) bool { 
	derefServerMap := *sh.serverMap
	if sh.scheme == "least-connections" {
		server_1 := derefServerMap[sh.items[i]].(*redis.Server)
		server_2 := derefServerMap[sh.items[j]].(*redis.Server)
		return server_1.GetActiveConnections() < server_2.GetActiveConnections()
	}
	return false
}

func (sh *ServerHeap) Swap(i, j int) { sh.items[i], sh.items[j] = sh.items[j], sh.items[i] }

func (sh *ServerHeap) Push(x string) {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	sh.items = append(sh.items, x)
}

func (sh *ServerHeap) Pop() string {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	n := len(sh.items)
	x := sh.items[n-1]
	sh.items = sh.items[0 : n-1]
	return x
}

