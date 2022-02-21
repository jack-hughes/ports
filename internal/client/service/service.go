package service

import (
	"context"
	"github.com/jack-hughes/ports/pkg/apis/ports"
	"google.golang.org/grpc"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -source=./service.go -package=mocks -destination=../../../test/mocks/ports_client_mocks.go
type PortClient interface {
	Update(ctx context.Context, req *ports.Request) error
	CloseAndRecv() (*ports.Response, error)
}

type Ports struct {
	client ports.PortsClient
	stream ports.Ports_UpdateClient
}

func NewPortsClient(ctx context.Context, conn grpc.ClientConnInterface) (Ports, error) {
	c := ports.NewPortsClient(conn)

	s, err := c.Update(ctx)
	if err != nil {
		return Ports{}, err
	}

	return Ports{
		client: c,
		stream: s,
	}, nil
}

func (p Ports) Update(ctx context.Context, req *ports.Request) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	if err := p.stream.Send(req); err != nil {
		return err
	}

	return nil
}

func (p Ports) CloseAndRecv() (*ports.Response, error) {
	return p.stream.CloseAndRecv()
}
