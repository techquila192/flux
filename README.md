# Flux
A simple, concurrent HTTP load balancer written in Golang and uses a Redis datastore. Currently supports two schemes for request distribution :
* Round robin (static)
* Least active connections (dynamic)

If the current candidate for request forwarding is not available, then the load balancer will execute a certain number of retries before marking the current server as inactive. The priority is given to request fulfillment rather than decisive rejection, potentially causing overutilization of resources to concurrently fulfill all requests if servers are unavailable for extended periods of time.

Health checks are performed at configurable intervals to periodically detect dead nodes and pick up newly alive nodes. 
Flux also supports dynamic addition and removal of new servers via HTTP requests to :
* **/add-server** 
* **/remove-server**

With 2 query parameters:
* *server* - Address of the server
* *code* - A secret configured in the config file 

## Usage

1. ### To run a Docker image 
* Clone the repository and navigate to root directory
* Create a env file with SECRET as key. This will be used to verify requests for dynamic addition and removal of servers.
* Configure the config file (**/config/config.json**) according to specifications.
* Build the docker image  
`docker build -t <image-name> .`
* Run the image in a container

2. ### To run the Go binary
* Clone the repository and navigate to root directory
* Ensure the Redis server is running (locally for minimal latency)
* Create a env file with SECRET as key. This will be used to verify requests for dynamic addition and removal of servers.
* Configure the config file (**/config/config.json**) according to specifications.
* Run `go mod download` to download dependencies
* Run `go run main.go` to start the application or alternatively, build the compiled binary and execute it.




