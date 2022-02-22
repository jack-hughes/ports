package options

import "flag"

// Options provides configuration options.
type Options struct {
	// gRPC Server options
	GRPCServer string
	GRPCPort   string

	// HTTP Server options
	HTTPServer string
	HTTPPort   string

	FilePath string
	LogLevel int
}

// DefaultOptions for the server
var DefaultOptions = Options{
	GRPCServer: "localhost",
	GRPCPort:   "50085",

	HTTPServer: "localhost",
	HTTPPort:   "8181",

	FilePath: "test/testdata/ports.json",

	// Default to debug logs
	LogLevel: -1,
}

// FillOptionsUsingFlags fill options from the command line / default options.
func (o *Options) FillOptionsUsingFlags(flags *flag.FlagSet) {
	flags.StringVar(&o.GRPCServer, "grpc-server", DefaultOptions.GRPCServer, "the hostname of the remote gRPC storage server")
	flags.StringVar(&o.GRPCPort, "grpc-port", DefaultOptions.GRPCPort, "the port of the remote gRPC storage server")

	flags.StringVar(&o.HTTPServer, "http-server", DefaultOptions.HTTPServer, "host that the gRPC server should use on the local interface")
	flags.StringVar(&o.HTTPPort, "http-port", DefaultOptions.HTTPPort, "port that the gRPC server should use on the local interface")

	flags.StringVar(&o.FilePath, "filepath", DefaultOptions.FilePath, "filepath for the document to load")
	flags.IntVar(&o.LogLevel, "log-level", DefaultOptions.LogLevel, "the log level for the application")
}
