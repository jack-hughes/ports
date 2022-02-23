package types

import "github.com/jack-hughes/ports/pkg/apis/ports"

// InMemStore is a map of Ports with the key as the port ID
type InMemStore struct {
	Ports map[string]Port
}

// PortStream to receive ports and errors on
type PortStream struct {
	ID    string
	Port  Port
	Error error
}

// Port a representation of port details
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

// Clone turns a gRPC specific type into an internal port type
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

// ToTransit turns an internal type into a gRPC specific port type
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
