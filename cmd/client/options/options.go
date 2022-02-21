package options

import "flag"

// Options provides command line configuration options for the client
type Options struct {
	// gRPC Server options
	GRPCServer string
	GRPCPort   string
}

// DefaultOptions for the client to operate
var DefaultOptions = Options{
	GRPCServer: "localhost",
	GRPCPort:   "50085",
}

// FillOptionsUsingFlags fill options from the command line / default options.
func (o *Options) FillOptionsUsingFlags(flags *flag.FlagSet) {
	flags.StringVar(&o.GRPCServer, "grpc-server", DefaultOptions.GRPCServer, "the hostname of the remote gRPC storage server")
	flags.StringVar(&o.GRPCPort, "grpc-port", DefaultOptions.GRPCPort, "the port of the remote gRPC storage server")
}
