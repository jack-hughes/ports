.PHONY: test build proto

race:
	mkdir -p artifacts
	go test -race -short -cover -count 1 -coverprofile=artifacts/coverage.txt -covermode=atomic ./...

generate:
	go mod tidy
	go mod vendor
	go generate ./...
	go fmt ./...
	go vet ./...

proto:
	protoc --go_out="pkg" --go-grpc_out="pkg" \
		--go_opt=paths=source_relative --go-grpc_opt=paths=source_relative \
		apis/ports/ports.proto

run-client:
	go run cmd/client/main.go

run-server:
	go run cmd/server/main.go

build: generate
	mkdir -p bin
	go build -o bin/client/client ./cmd/client
	go build -o bin/server/server ./cmd/server
