package redis


import (
    "context"
    "github.com/go-redis/redis/v8"
)

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sync"
)

type Server struct {
	host *url.URL 
	activeConnections int 
	isAlive bool
	mutex sync.Mutex
}
func (s *Server) setHost(host string){
	URLhost, err := url.Parse(host)
	s.host = URLhost
	if err!=nil{
		fmt.Println("Error parsing url",err)
		return
	}

}

func (s *Server) getHost() *url.URL{
	return s.host
}

func (s *Server) getIsAlive() bool {
	return s.isAlive
}

func (s *Server) setAliveState(flag bool){
	s.isAlive=flag 
}

func (s *Server) getActiveConnections() int {
	return s.activeConnections;
}

func (s *Server) incrementConnections() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.activeConnections++
}

func (s *Server) decrementConnections() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.activeConnections--
}


var	Servers map[*Server]bool 
var client 	*redis.Client //read only


func Connect(redis_host *url.URL, get_dump bool, dump_location ...string) bool{
	redis_host.Scheme = ""
	Servers = make(map[*Server]bool)
	client = redis.NewClient(&redis.Options{
        Addr:     redis_host.String(), // Redis server address
        Password: "",               // No password
        DB:       0,                // Default DB
    })
	ctx:= context.Background()
	_, err := client.Ping(ctx).Result()
    if err != nil {
        return false
    } else {

		if get_dump  {
		// Set the directory of the dump file
		_, err := client.ConfigSet(ctx, "dir", dump_location[0]).Result()
		if err != nil {
			fmt.Println("Error setting directory:", err)
			return false
		}

		// Set the filename of the dump file
		_, err = client.ConfigSet(ctx, "dbfilename", dump_location[1]).Result()
		if err != nil {
			fmt.Println("Error setting dbfilename:", err)
			return false
		}
	}
		return true
	}

}

func InitServerList(init_servers *[]string){
	serverSet := GetSetfromDump()

}

func GetSetfromDump() {
	serverSet, err := client.SMembers(context.Background(), "Servers").Result()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}





