package main

import (
	"context"
	"flag"
	"github.com/jack-hughes/ports/cmd/client/options"
	"github.com/jack-hughes/ports/internal/client/service"
	"github.com/jack-hughes/ports/internal/client/stream"
	"github.com/jack-hughes/ports/pkg/apis/ports"
	"google.golang.org/grpc"
	"log"
	"net"
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
	c, err := service.NewPortsClient(ctx, conn)
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for data := range s.Watch() {
			if data.Error != nil {
				log.Fatal("failure to read file - exiting")
			}

			err := c.Update(ctx, &ports.Request{
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
}
