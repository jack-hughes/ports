package service

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	types "github.com/jack-hughes/ports/internal"
	"github.com/jack-hughes/ports/internal/server/storage"
	"github.com/jack-hughes/ports/pkg/apis/ports"
	"io"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -package mocks -source=../../../pkg/apis/ports/ports_grpc.pb.go -destination=../../../test/mocks/ports_server_mocks.go -build_flags=-mod=mod
type PortsServer struct {
	ports.UnimplementedPortsServer

	store storage.Storage
}

// New instantiates a new PortsServer
func New(s storage.Storage) *PortsServer {
	return &PortsServer{
		store: s,
	}
}

// Update port service function
func (s *PortsServer) Update(stream ports.Ports_UpdateServer) error {
	for {
		port, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&empty.Empty{})
		}
		if err != nil {
			return err
		}

		s.store.Update(types.Clone(port))
	}
}

func (s *PortsServer) Get(ctx context.Context, p *ports.GetPortRequest) (*ports.Port, error) {
	return types.ToTransit(s.store.Get(p.ID)), nil
}

func (s *PortsServer) List(e *empty.Empty, srv ports.Ports_ListServer) error {
	l := s.store.List()
	for _,v := range l {
		if err := srv.Send(types.ToTransit(v)); err != nil {
			return err
		}

	}

	return nil
}
