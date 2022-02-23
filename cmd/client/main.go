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
	"os"
	"os/signal"
	"syscall"
	"time"
)

// AppName constant provided to logger
const AppName = "ports-client-api"

func main() {
	opts := options.DefaultOptions
	opts.FillOptionsUsingFlags(flag.CommandLine)
	flag.Parse()
	ctx := context.TODO()

	log := logger.NewZapLogger(AppName, zapcore.Level(opts.LogLevel))
	log.Debug("booting...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		sig := <-sigChan
		log.Info(fmt.Sprintf("system call: %v", sig))
		cancel()
	}()

	conn, err := grpc.Dial(net.JoinHostPort(opts.GRPCServer, opts.GRPCPort), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatal("failed to connect to server: %v", zap.Error(err))
	}

	// Set up a JSON stream
	s := stream.NewJSONStream(log)
	c, err := service.NewService(ctx, conn, log)
	if err != nil {
		log.Fatal("failed to create service: %v", zap.Error(err))
	}

	// Concurrently wait for data and send it via an update on a client -> server stream when it arrives in the channel
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

	// Start reading the JSON file
	s.Start(opts.FilePath)

	// Start the HTTP server
	if err := serve(ctx, opts, log, c); err != nil {
		log.Fatal("failed to serve: %v", zap.Error(err))
	}
}

func serve(ctx context.Context, opts options.Options, log *zap.Logger, c service.Service) (err error) {
	r := mux.NewRouter()
	r.HandleFunc("/ports", handlers.List(ctx, c, log)).Methods(http.MethodGet)
	r.HandleFunc("/ports/{port_id}", handlers.Get(ctx, c, log)).Methods(http.MethodGet)

	addr := net.JoinHostPort(opts.HTTPServer, opts.HTTPPort)
	srv := &http.Server{Addr: addr, Handler: r}

	go func() {
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to start http server: %v", zap.Error(err))
		}
	}()

	log.Info(fmt.Sprintf("serving on: %s", addr))
	<-ctx.Done()
	log.Info("server shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	if err = srv.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown failed: %w", zap.Error(err))
	}

	log.Info("shutdown success")

	if err == http.ErrServerClosed {
		err = nil
	}

	return
}
