package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jack-hughes/ports/cmd/client/options"
	"github.com/jack-hughes/ports/internal/client/handlers"
	"github.com/jack-hughes/ports/internal/client/service"
	"github.com/jack-hughes/ports/internal/client/stream"
	"github.com/jack-hughes/ports/internal/logger"
	"github.com/jack-hughes/ports/pkg/apis/ports"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"net"
	"net/http"
)

const AppName = "ports-client"

func main() {
	opts := options.DefaultOptions
	opts.FillOptionsUsingFlags(flag.CommandLine)
	flag.Parse()
	ctx := context.TODO()


	log := logger.NewZapLogger(AppName, zapcore.Level(opts.LogLevel))
	conn, err := grpc.Dial(net.JoinHostPort(opts.GRPCServer, opts.GRPCPort), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal("failed to connect to server: %v", zap.Error(err))
	}

	s := stream.NewJSONStream(log)
	c, err := service.NewService(ctx, conn, log)
	if err != nil {
		log.Fatal("failed to create service: %v", zap.Error(err))
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
				log.Fatal("failed to send to server: %v", zap.Error(err))
			}
		}

		_, err := c.CloseAndRecv()
		if err != nil {
			log.Fatal("stream responded with error: %v", zap.Error(err))
		}
	}()

	s.Start(opts.FilePath)

	r := mux.NewRouter()
	r.HandleFunc("/ports", handlers.List(ctx, c, log)).Methods(http.MethodGet)
	r.HandleFunc("/ports/{port_id}", handlers.Get(ctx, c, log)).Methods(http.MethodGet)
	addr := net.JoinHostPort(opts.HTTPServer, opts.HTTPPort)
	log.Info(fmt.Sprintf("serving on: %s", addr))
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal("failed to start http server: %v", zap.Error(err))
	}
}
