package types

import "github.com/jack-hughes/ports/pkg/apis/ports"

type InMemStore struct {
	Ports map[string]Port
}

type PortStream struct {
	ID    string
	Port  Port
	Error error
}

type Error struct {
	Error string `json:"error"`
}

type Port struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	City        string    `json:"city"`
	Country     string    `json:"country"`
	Alias       []string  `json:"alias"`
	Regions     []string  `json:"regions"`
	Coordinates []float32 `json:"coordinates"`
	Province    string    `json:"province"`
	Timezone    string    `json:"timezone"`
	Unlocs      []string  `json:"unlocs"`
	Code        string    `json:"code"`
}

func Clone(req *ports.Port) Port {
	return Port{
		ID:          req.ID,
		Name:        req.Name,
		City:        req.City,
		Country:     req.Country,
		Alias:       req.Alias,
		Regions:     req.Regions,
		Coordinates: req.Coordinates,
		Province:    req.Province,
		Timezone:    req.Timezone,
		Unlocs:      req.Unlocs,
		Code:        req.Code,
	}
}

func ToTransit(req Port) *ports.Port {
	return &ports.Port{
		ID:          req.ID,
		Name:        req.Name,
		City:        req.City,
		Country:     req.Country,
		Alias:       req.Alias,
		Regions:     req.Regions,
		Coordinates: req.Coordinates,
		Province:    req.Province,
		Timezone:    req.Timezone,
		Unlocs:      req.Unlocs,
		Code:        req.Code,
	}
}
