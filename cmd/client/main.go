package main

import (
	"context"
	"flag"
	"github.com/gorilla/mux"
	"github.com/jack-hughes/ports/cmd/client/options"
	"github.com/jack-hughes/ports/internal/client/handlers"
	"github.com/jack-hughes/ports/internal/client/service"
	"github.com/jack-hughes/ports/internal/client/stream"
	"github.com/jack-hughes/ports/pkg/apis/ports"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
)

func main() {
	opts := options.DefaultOptions
	opts.FillOptionsUsingFlags(flag.CommandLine)
	flag.Parse()
	ctx := context.TODO()

	conn, err := grpc.Dial(net.JoinHostPort(opts.GRPCServer, opts.GRPCPort), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalln(err)
	}

	s := stream.NewJSONStream()
	c, err := service.NewService(ctx, conn)
	if err != nil {
		log.Fatal("failed to create service")
	}
	go func() {
		for data := range s.Watch() {
			if data.Error != nil {
				log.Fatal("failure to read file - exiting")
			}

			err := c.Update(ctx, &ports.Port{
				ID:          data.ID,
				Name:        data.Port.Name,
				City:        data.Port.City,
				Country:     data.Port.Country,
				Alias:       data.Port.Alias,
				Regions:     data.Port.Regions,
				Coordinates: data.Port.Coordinates,
				Province:    data.Port.Province,
				Timezone:    data.Port.Timezone,
				Unlocs:      data.Port.Unlocs,
				Code:        data.Port.Code,
			})
			if err != nil {
				log.Fatalf("failed to send to server: %v", err)
			}
		}

		_, err := c.CloseAndRecv()
		if err != nil {
			log.Fatalf("stream responded with error: %v", err)
		}
	}()

	s.Start("test/testdata/ports.json")

	r := mux.NewRouter()
	r.HandleFunc("/ports", handlers.List(ctx, c)).Methods(http.MethodGet)
	r.HandleFunc("/ports/{port_id}", handlers.Get(ctx, c)).Methods(http.MethodGet)

	log.Println("serving...")
	log.Fatal(http.ListenAndServe(net.JoinHostPort(opts.HTTPServer, opts.HTTPPort), r))
}
