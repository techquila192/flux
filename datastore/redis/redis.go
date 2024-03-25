package redis


import (
    "context"
    "github.com/go-redis/redis/v8"
)

import (
	"fmt"
	"sync"
	"strings"
)


type Server struct {
	Host string
	ActiveConnections int 
	IsAlive bool
	Mutex sync.Mutex
}




func (s *Server) SetHost(host string){
	s.Host = strings.Trim(host," ")

}


func (s *Server) GetHost() string{
	return s.Host
}

func (s *Server) GetIsAlive() bool {
	return s.IsAlive
}

func (s *Server) SetAliveState(flag bool){
	s.IsAlive=flag 
}

func (s *Server) GetActiveConnections() int {
	return s.ActiveConnections;
}

func (s *Server) IncrementConnections() {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.ActiveConnections++
}

func (s *Server) DecrementConnections() {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.ActiveConnections--
}



type Redis struct {
	Ctx context.Context
	client 	*redis.Client
	Mu sync.RWMutex  //mutex for server list
	Servers map[string] *Server
	Get_Dump bool
	
}


//configure redis to run with necessary dir to save and load dump
func (r *Redis) Connect(redis_host string, get_dump bool, dump_interval string) bool{
	redis_host = strings.Trim(redis_host," ")
	r.Get_Dump = get_dump
	r.Servers = make(map[string]*Server)
	r.client = redis.NewClient(&redis.Options{
        Addr:     redis_host, 		// Redis server address
        Password: "",               // No password
        DB:       0,                // Default DB
    })
	r.Ctx = context.Background()
	_, err := r.client.Ping(r.Ctx).Result()
    if err != nil {
		fmt.Println(err)
        return false
    } else {
		r.client.ConfigSet(r.Ctx, "save", dump_interval+" 1") // set snapshot interval and set to save after change
		return true
	}

}

func (r *Redis) InitServerList(init_servers *[]string) (*map[string]*Server){
	if !r.Get_Dump {
		_, err := r.client.Del(r.Ctx,"Servers").Result()
		if err != nil {
			fmt.Println("Error deleting set:", err)
		}
	} else {
		load_status := r.SyncWithRedisSet()
		if !load_status {
			fmt.Println("Could not sync with redis")
		}
		return &r.Servers
	}
	
	
	for _, serv_ip := range (*init_servers){
		server_instance := Server{}
		server_instance.SetHost(serv_ip)
		r.Servers[server_instance.GetHost()]=&server_instance
		//add to redis 
		_, err := r.client.SAdd(r.Ctx, "Servers", server_instance.GetHost()).Result()
		if err != nil {
			fmt.Print("Error while adding member",err)
			panic(err)
		}
	}

	return &r.Servers

}

//will reset all object
func (r *Redis) SyncWithRedisSet() bool{
	
	redisSet, err := r.client.SMembers(r.Ctx, "Servers").Result() //returns array of string
	
	if err != nil {
		fmt.Println("Error in getting servers from dump", err)
		return false
	}
	fmt.Println("sync: ",redisSet)
	r.Mu.RLock()
	defer r.Mu.RUnlock()
	for _,host_name := range redisSet{
		server_instance := Server{}
		server_instance.SetHost(host_name)
		r.Servers[server_instance.GetHost()]=&server_instance
		
	}
	
	return true
}

func (r *Redis) AddMember(server_host string){
	server_instance := Server{}
	server_instance.SetHost(strings.Trim(server_host, " "))
	_, present := r.Servers[server_instance.GetHost()]
	if present {
		return
	}
	r.Mu.Lock() //write lock
	defer r.Mu.Unlock()
	r.Servers[server_instance.GetHost()]=&server_instance
	// modify redis set
	_, err := r.client.SAdd(r.Ctx, "Servers", server_instance.GetHost()).Result()
	if err != nil {
		fmt.Println("Error while adding member",err)
		panic(err)
	}
	

}

func (r *Redis) RemoveMember(server string){
	r.Mu.Lock()
	defer r.Mu.Unlock()
	server = strings.Trim(server, " ")
	_, present := r.Servers[server]
	if !present {
		return
	}
	result := r.client.SRem(r.Ctx, "Servers", server)
    if err := result.Err(); err != nil {
        fmt.Print("Error while deleting member")
    }
	delete(r.Servers,server)
	

}

func (r *Redis) GetServers() (*map[string]*Server){ 

	return &r.Servers
	
}

func (r *Redis) GetClient() (*redis.Client){

	return r.client
}

