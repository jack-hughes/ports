package main

import (
	"flag"
	"fmt"
	"github.com/jack-hughes/ports/cmd/server/options"
	"github.com/jack-hughes/ports/internal/server/service"
	"github.com/jack-hughes/ports/pkg/apis/ports"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	opts := options.DefaultOptions
	opts.FillOptionsUsingFlags(flag.CommandLine)
	flag.Parse()

	addr := net.JoinHostPort(opts.GRPCServer, opts.GRPCPort)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal("failed to listen")
	}

	srv := grpc.NewServer()
	svc := service.New()
	ports.RegisterPortsServer(srv, svc)

	log.Println(fmt.Sprintf("listening on %s", addr))
	if err := srv.Serve(lis); err != nil {
		log.Fatal("failed to serve")
	}
}
