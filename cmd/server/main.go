package main

import (
	"flag"
	"fmt"
	"github.com/jack-hughes/ports/cmd/server/options"
	"github.com/jack-hughes/ports/internal/logger"
	"github.com/jack-hughes/ports/internal/server/service"
	"github.com/jack-hughes/ports/internal/server/storage"
	"github.com/jack-hughes/ports/pkg/apis/ports"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"net"
)

const AppName = "port-domain-service"

func main() {
	opts := options.DefaultOptions
	opts.FillOptionsUsingFlags(flag.CommandLine)
	flag.Parse()

	log := logger.NewZapLogger(AppName, zapcore.Level(opts.LogLevel))
	log.Debug("booting...")
	addr := net.JoinHostPort(opts.GRPCServer, opts.GRPCPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("failed to listen", zap.Error(err))
	}

	srv := grpc.NewServer()
	store := storage.NewStorage(log)
	svc := service.New(store, log)
	ports.RegisterPortsServer(srv, svc)

	log.Debug(fmt.Sprintf("listening on %s", addr))
	if err := srv.Serve(lis); err != nil {
		log.Fatal("failed to serve")
	}
}
