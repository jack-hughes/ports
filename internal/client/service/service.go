package service

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	types "github.com/jack-hughes/ports/internal"
	"github.com/jack-hughes/ports/pkg/apis/ports"
	"google.golang.org/grpc"
	"io"
	"log"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -source=./service.go -package=mocks -destination=../../../test/mocks/ports_client_mocks.go
type PortClient interface {
	Update(ctx context.Context, port types.Port) error
    Get(ctx context.Context, id string) (types.Port, error)
	List(ctx context.Context) ([]types.Port, error)
	CloseAndRecv() (*empty.Empty, error)
}

type Service struct {
	client ports.PortsClient
	stream ports.Ports_UpdateClient
}

func NewService(ctx context.Context, conn grpc.ClientConnInterface) (Service, error) {
	c := ports.NewPortsClient(conn)

	s, err := c.Update(ctx)
	if err != nil {
		return Service{}, err
	}

	return Service{
		client: c,
		stream: s,
	}, nil
}

func (s Service) Update(ctx context.Context, port *ports.Port) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := s.stream.Send(port); err != nil {
		return err
	}

	return nil
}

func (s Service) Get(ctx context.Context, id string) (types.Port, error) {
	resp, err := s.client.Get(ctx, &ports.GetPortRequest{ID: id})
	if err != nil {
		return types.Port{}, err
	}

	return types.Clone(resp), nil
}

func (s Service) List(ctx context.Context) ([]types.Port, error) {
	stream, err := s.client.List(ctx, &empty.Empty{})
	if err != nil {
		return nil, err
	}

	var pSlice []types.Port
	done := make(chan bool)
	go func() {
		for {
			log.Println("RECEIVING")
			p, err := stream.Recv()
			if err == io.EOF {
				done <- true
				return
			}
			if err != nil {
				return
			}

			pSlice = append(pSlice, types.Clone(p))
		}
	}()

	<-done
	return pSlice, nil
}

func (s Service) CloseAndRecv() (*empty.Empty, error) {
	return s.stream.CloseAndRecv()
}
