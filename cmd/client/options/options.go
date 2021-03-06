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

	// The path of the file to decode and send for storage
	FilePath string

	// Log levels in line with Zap
	// See here: https://github.com/uber-go/zap/blob/master/level.go
	LogLevel int
}

// DefaultOptions for the server
var DefaultOptions = Options{
	GRPCServer: "ports-domain-service",
	GRPCPort:   "50085",

	HTTPServer: "",
	HTTPPort:   "8181",

	FilePath: "/test/ports.json",

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
