package storage

import (
	types "github.com/jack-hughes/ports/internal"
	"log"
)

type Storage interface {
	Update(port types.Port) types.Port
	Get(portID string) types.Port
	List() []types.Port
}

type Store struct {
	db types.InMemStore
}

func NewStorage() Storage {
	return Store{db: types.InMemStore{
		Ports: make(map[string]types.Port),
	}}
}

func (s Store) Update(port types.Port) types.Port {
	s.db.Ports[port.ID] = port
	return s.db.Ports[port.ID]
}

func (s Store) Get(portID string) types.Port {
	return s.db.Ports[portID]
}

func (s Store) List() []types.Port {
	var list []types.Port
	for _, v := range s.db.Ports {
		log.Println("hitting")
		list = append(list, v)
	}
	return list
}