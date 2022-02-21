package options

import "flag"

// Options provides configuration options.
type Options struct {
	// gRPC Server options
	GRPCServer string
	GRPCPort   string
}

// DefaultOptions for the server
var DefaultOptions = Options{
	GRPCServer: "localhost",
	GRPCPort:   "50085",
}

// FillOptionsUsingFlags Fill options from the command line / default options.
func (o *Options) FillOptionsUsingFlags(flags *flag.FlagSet) {
	flags.StringVar(&o.GRPCServer, "grpc-server", DefaultOptions.GRPCServer, "host that the gRPC server should use on the local interface")
	flags.StringVar(&o.GRPCPort, "grpc-port", DefaultOptions.GRPCPort, "port that the gRPC server should use on the local interface")
}
