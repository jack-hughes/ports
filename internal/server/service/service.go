package service

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	types "github.com/jack-hughes/ports/internal"
	"github.com/jack-hughes/ports/internal/server/storage"
	"github.com/jack-hughes/ports/pkg/apis/ports"
	"go.uber.org/zap"
	"io"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -package mocks -source=../../../pkg/apis/ports/ports_grpc.pb.go -destination=../../../test/mocks/ports_server_mocks.go -build_flags=-mod=mod
type PortsServer struct {
	ports.UnimplementedPortsServer

	store storage.Storage
	log   *zap.Logger
}

// New instantiates a new PortsServer
func New(s storage.Storage, log *zap.Logger) *PortsServer {
	return &PortsServer{
		store: s,
		log:   log.With(zap.String("component", "service")),
	}
}

// Update port service function
func (s *PortsServer) Update(stream ports.Ports_UpdateServer) error {
	for {
		s.log.Debug("attempting receiving a port update")
		port, err := stream.Recv()
		if err == io.EOF {
			s.log.Debug("end of stream, nothing to receive")
			return stream.SendAndClose(&empty.Empty{})
		}
		if err != nil {
			s.log.Error("failed to receive item: %v", zap.Error(err))
			return err
		}

		s.log.Debug(fmt.Sprintf("updating store for port: %v", port.ID))
		s.store.Update(types.Clone(port))
	}
}

func (s *PortsServer) Get(ctx context.Context, p *ports.GetPortRequest) (*ports.Port, error) {
	s.log.Debug(fmt.Sprintf("attempting to retrieve port entry for id: %v", p.ID))

	port, err := s.store.Get(p.ID)
	if err != nil {
		s.log.Error("failed to get item: %v", zap.Error(err))
		return &ports.Port{}, err
	}

	s.log.Debug(fmt.Sprintf("returning port entry for id: %v", p.ID))
	return types.ToTransit(port), nil
}

func (s *PortsServer) List(e *empty.Empty, srv ports.Ports_ListServer) error {
	s.log.Debug("attempting to list all ports")
	l := s.store.List()
	for _, v := range l {
		s.log.Debug(fmt.Sprintf("attempting to send port on stream: %v", v.ID))
		if err := srv.Send(types.ToTransit(v)); err != nil {
			s.log.Error("failed to send item: %v", zap.Error(err))
			return err
		}

	}
	s.log.Debug("completed listing ports")

	return nil
}
