package service

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	types "github.com/jack-hughes/ports/internal"
	"github.com/jack-hughes/ports/pkg/apis/ports"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"io"
)

// Service contains objects required by the service to opertate. In this case,
// the client to interact with and the stream to send updates on
//go:generate go run -mod=mod github.com/golang/mock/mockgen -package mocks -source=../../../pkg/apis/ports/ports_grpc.pb.go -destination=../../../test/mocks/grpc_mocks.go -build_flags=-mod=mod
type Service struct {
	client ports.PortsClient
	stream ports.Ports_UpdateClient
	log    *zap.Logger
}

// NewService instantiates a new service with a gRPC client client connection and logger
func NewService(ctx context.Context, conn grpc.ClientConnInterface, log *zap.Logger) (Service, error) {
	c := ports.NewPortsClient(conn)

	s, err := c.Update(ctx)
	if err != nil {
		return Service{}, err
	}

	return Service{
		client: c,
		stream: s,
		log:    log.With(zap.String("component", "service")),
	}, nil
}

// Update a record on the remote server, if it doesn't exist, create it
func (s Service) Update(ctx context.Context, port *ports.Port) error {
	s.log.Debug(fmt.Sprintf("incoming update with port ID: %v", port.ID))

	if err := s.stream.Send(port); err != nil {
		s.log.Error("failed to send to stream: %v", zap.Error(err))
		return err
	}

	return nil
}

// Get an individual record from the remote server
func (s Service) Get(ctx context.Context, id string) (types.Port, error) {
	s.log.Debug(fmt.Sprintf("attempting to get port for ID: %v", id))
	resp, err := s.client.Get(ctx, &ports.GetPortRequest{ID: id})
	if err != nil {
		s.log.Error("failed to get port: %v", zap.Error(err))
		return types.Port{}, err
	}

	return types.Clone(resp), nil
}

// List all records in storage
func (s Service) List(ctx context.Context) ([]types.Port, error) {
	s.log.Debug("attempting to list all available ports")
	stream, err := s.client.List(ctx, &empty.Empty{})
	if err != nil {
		s.log.Error("failed to list ports: %v", zap.Error(err))
		return nil, err
	}

	var pSlice []types.Port
	done := make(chan bool)
	go func() {
		for {
			s.log.Debug("receiving ports")
			p, err := stream.Recv()
			if err == io.EOF {
				s.log.Debug("end of ports stream")
				done <- true
				return
			}
			if err != nil {
				s.log.Error("failed to receive port: %v", zap.Error(err))
			}

			pSlice = append(pSlice, types.Clone(p))
		}
	}()

	<-done
	return pSlice, nil
}

// CloseAndRecv closes the stream
func (s Service) CloseAndRecv() (*empty.Empty, error) {
	s.log.Debug("closing stream")
	return s.stream.CloseAndRecv()
}
