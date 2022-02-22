package storage

import (
	"fmt"
	types "github.com/jack-hughes/ports/internal"
	"go.uber.org/zap"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -package mocks -source=./storage.go -destination=../../../test/mocks/storage_mocks.go -build_flags=-mod=mod
type Storage interface {
	Update(port types.Port) types.Port
	Get(portID string) (types.Port, error)
	List() []types.Port
}

type Store struct {
	db  types.InMemStore
	log *zap.Logger
}

func NewStorage(log *zap.Logger) Storage {
	return Store{
		db:  types.InMemStore{Ports: make(map[string]types.Port)},
		log: log.With(zap.String("component", "storage")),
	}
}

func (s Store) Update(port types.Port) types.Port {
	s.db.Ports[port.ID] = port
	s.log.Debug(fmt.Sprintf("port updated: %v", port.ID))
	return s.db.Ports[port.ID]
}

func (s Store) Get(portID string) (types.Port, error) {
	for _, v := range s.db.Ports {
		if v.ID == portID {
			return v, nil
		}
	}
	return types.Port{}, fmt.Errorf("could not find port with id: %v", portID)
}

func (s Store) List() []types.Port {
	var list []types.Port
	for _, v := range s.db.Ports {
		list = append(list, v)
	}
	return list
}
