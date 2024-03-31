# Request Counter

## Requirements

### Developer Requirements

Deliver a "RequestCount" server, written in Go. This server...
* must count the number of requests handled by the server instance
* must count the total number of requests handled by all servers that are running
* must serve these numbers through HTTP, with text/plain response content

Dockerize the "RequestCount" server
* tip 1: build the binary with CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "main_name" main_name.go
* tip 2: you also can use one of the golang Docker base images
  Deploy multiple instances of the "RequestCount" server with Docker Compose (replicas: 3)

#### Example flow:

```
$ curl http://localhost:8083
You are talking to instance host2:8083.
This is request 1 to this instance and request 1 to the cluster.
```

```
$ curl http://localhost:8083
You are talking to instance host1:8083.
This is request 1 to this instance and request 2 to the cluster.
```

```
$ curl http://localhost:8083
You are talking to instance host2:8083.
This is request 2 to this instance and request 3 to the cluster.
```

### DevOps Requirements

Make deployment templates for the "RequestCount" server. These deployment templates...
* must be Helm 3 charts
* must not include unnecessary templates
* must be self contained: service(s) for synchronizing state between the "RequestCount" servers must be deployed too

## Info

* The project layout follows https://github.com/golang-standards/project-layout.
    * `build/package` : Dockerfile.
    * `cmd/requestcounter` : entrypoint of the application.
    * `deployments/docker-compose` : docker-compose configuration.
    * `deployments/heml` : helm chart.
    * `internal/app/requestcounter` : business logic and infrastructure independent components.
    * `internal/infra` : infrastructure dependent components.
* Application reads configuration from environment variables using `kelseyhightower/envconfig`.
* `atomic.Int64` is used to count the number of requests handled by the server instance.
* Redis `INCR` command is used to count the total number of requests handled by all servers that are running.
* There is no persistence configuration for Redis.
* `redis/go-redis` is used as Redis client.
* `stretchr/testify` is used for testing utilities.
* `orlangure/gnomock` is used for running containers during testing.

## Running Tests

To run tests, run the following command. It requires Docker.

```shell
make test
```

## Run in Docker Compose

To run locally, run the following command. It will expose application `localhost:8083`.

```shell
make docker-compose
```

## Run in Kubernetes

To run locally, run the following command. It will expose application `localhost:8083`.

```shell
make helm
```