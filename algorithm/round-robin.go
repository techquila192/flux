package algorithm

import (
	"flux/datastore/redis"
	"sync"
)



func Round_robin(servers *map[string]*redis.Server, servers_mu *sync.RWMutex, vis *map[string]bool, vis_mu *sync.Mutex) (*redis.Server) {
	 servers_mu.RLock()
	 vis_mu.Lock()
	 defer servers_mu.RUnlock()
	 defer vis_mu.Unlock()
	 vis_map := *vis
	 var available int = 0
	 for name, server := range *servers {
		if !server.GetIsAlive(){
			continue
		}
		_, check := vis_map[name]
		if !check {
			available++
			break
		}
	}
	if available == 0{
		*vis = make(map[string]bool)
		vis_map = *vis
	} 
	 for name, server := range *servers {
		if !server.GetIsAlive(){
			continue
		}
		_, check := vis_map[name]
		if !check {
			vis_map[name] = true
			return server
		}
	 }
	 return nil

}