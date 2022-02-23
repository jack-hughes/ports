# Ports
[![codecov](https://codecov.io/gh/jack-hughes/ports/branch/main/graph/badge.svg?token=SCHIWMXR8K)](https://codecov.io/gh/jack-hughes/ports)
[![CI](https://github.com/jack-hughes/ports/actions/workflows/ci.yml/badge.svg)](https://github.com/jack-hughes/ports/actions/workflows/ci.yml)

## About
Ports is made up of two services. The first, the `port-client-service`, is responsible for parsing a JSON file of pre-determined structure (containing, coincidentally, information on shipping ports). When the application parses this file it sends each chunk on a gRPC stream to be stored in memory within the second service, the `ports-domain-service`. Once this is completed, a HTTP server is exposed on the client service that allows users to query the ingested data held in the `ports-domain-service`.

## Installation 
### Requirements
- [Go](https://go.dev/) (built using version 1.17.2)
- [Docker Compose](https://docs.docker.com/compose/install/)
- [Protocol Buffers](https://developers.google.com/protocol-buffers/docs/gotutorial)

### Running
You can run both the `ports-client-service` and `port-domain-service` locally by running the corresponding make commands in two terminal sessions. In one, simply run `make run-server`. The server should be run first to ensure that the client can successfully connect when launched. When the server is up, you can run `make run-client` in the second terminal window.

Alternatively, you can run the application with `docker-compose` by running `make up`, which will build both containers and bring them both up. To check the logs of these containers, you can run `make logs`. To gracefully shut down the compose stack, run `make down`.

### Logging
By default, both applications start with debug level logging. When the applications are started, this is incredibly verbose. To turn this down, you can supply the `--log-level` flag to both client and server. Please follow Zap's [standards](https://github.com/uber-go/zap/blob/master/level.go) for log levels when setting this  value.
### Tooling
- `make integration`
  - Run integration tests when we have either the local or `docker-compose` stack running
- `make race`
  - Run unit tests with coverage across the codebase
- `make generate`
  - Run `go mod tidy`, `go mod vendor`, `go generate` for mocks, `go fmt` and `go vet`
- `make proto`
  - Generate the protocol buffers for the services
- `make build`
  - Generate go binaries for client and server, located in `bin/`
- `make up`
  - Bring up the docker compose stack
- `make logs`
  - Get logs out of the docker compose stack
- `make down`
    - Kill the docker compose stack
