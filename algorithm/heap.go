package algorithm

import (
	"flux/datastore/redis"
	"sync"
)


type ServerHeap struct{
	items []string
	Scheme string
	heap_mu sync.Mutex 
	mu sync.Mutex //for items 
	ServerMap *map[string] *redis.Server
}

func (sh *ServerHeap) Len() int { return len(sh.items) }

func (sh *ServerHeap) Less(i, j int) bool {
	derefServerMap := *sh.ServerMap
	if sh.Scheme == "least-connections" {
		
		server_1 := derefServerMap[sh.items[i]]
		server_2 := derefServerMap[sh.items[j]]
		if server_1 == nil && server_2 == nil {
			return true
		}
		if server_1 == nil{
			return false
		}
		if server_2 == nil{
			return true
		}
		return server_1.GetActiveConnections() < server_2.GetActiveConnections()
	}
	return false
}

func (sh *ServerHeap) Swap(i, j int) { sh.items[i], sh.items[j] = sh.items[j], sh.items[i] }

func (sh *ServerHeap) Push(x any) {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	sh.items = append(sh.items, x.(string))
}

func (sh *ServerHeap) Pop() any {
	sh.mu.Lock()
	defer sh.mu.Unlock()
	n := len(sh.items)
	x := sh.items[n-1]
	sh.items = sh.items[0 : n-1]
	return x
}

