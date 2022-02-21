package service

import (
	"github.com/jack-hughes/ports/pkg/apis/ports"
	"io"
	"log"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -package mocks -source=../../../pkg/apis/ports/ports_grpc.pb.go -destination=../../../test/mocks/ports_server_mocks.go -build_flags=-mod=mod
type PortsServer struct {
	ports.UnimplementedPortsServer
}

// New instantiates a new PortsServer
func New() *PortsServer {
	return &PortsServer{}
}

// Update port service function
func (s *PortsServer) Update(stream ports.Ports_UpdateServer) error {
	total := 0
	for {
		port, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&ports.Response{Todo: "todo"})
		}
		if err != nil {
			return err
		}
		total++
		log.Printf("total received: %v\n", total)
		log.Printf("port: %v\n", port)
	}
}
