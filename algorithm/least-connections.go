package algorithm

import (
	"container/heap"
	"flux/datastore/redis"
)

func Least_connections(sh *ServerHeap) (*redis.Server){
	sh.heap_mu.Lock()
	defer sh.heap_mu.Unlock()
	notAliveServers := make([]string,0)
	
	for sh.Len()!=0{
		server := heap.Pop(sh).(string)
		serverPtr := (*sh.ServerMap)[server]
		if serverPtr == nil{
			continue
		}
		if !serverPtr.GetIsAlive(){
			notAliveServers = append(notAliveServers, server)
		} else {
			//push ignored servers back
			for _, dead_server := range notAliveServers {
				heap.Push(sh, dead_server)
			}
			heap.Push(sh, server)
			return serverPtr
		}
	}
	
	return nil
}