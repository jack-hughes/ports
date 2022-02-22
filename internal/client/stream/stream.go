package stream

import (
	"encoding/json"
	"fmt"
	types "github.com/jack-hughes/ports/internal"
	"go.uber.org/zap"
	"os"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -package mocks -source=./stream.go -destination=../../../test/mocks/stream_mocks.go -build_flags=-mod=mod
// Streamer is the interface definition for actions on the PortStream
type Streamer interface {
	Watch() <-chan types.PortStream
	Start(path string)
}

// Stream contains the channel for PortStream objects
type Stream struct {
	stream chan types.PortStream
	log *zap.Logger
}

// NewJSONStream instantiates the new PortStream channel
func NewJSONStream(log *zap.Logger) Stream {
	return Stream{
		stream: make(chan types.PortStream),
		log: log.With(zap.String("component", "stream")),
	}
}

// Watch returns the contents of the stream
func (s Stream) Watch() <-chan types.PortStream {
	s.log.Debug("starting to watch")
	return s.stream
}

// Start opens a file and begins to decode the JSON object at the first delimiter. The decoder will step through each
// Port element, send it on the channel and exit when no more exist within the file
func (s Stream) Start(path string) {
	s.log.Debug(fmt.Sprintf("streaming file contents %v", path))
	defer close(s.stream)

	file, err := os.Open(path)
	if err != nil {
		s.stream <- types.PortStream{Error: fmt.Errorf("open file: %w", err)}
		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if _, err := decoder.Token(); err != nil {
		s.stream <- types.PortStream{Error: fmt.Errorf("decode opening delimiter: %w", err)}
		return
	}

	for decoder.More() {
		tkn, err := decoder.Token()
		if err != nil {
			s.stream <- types.PortStream{Error: fmt.Errorf("decode opening delimiter: %w", err)}
		}
		var port types.Port
		if err := decoder.Decode(&port); err != nil {
			s.stream <- types.PortStream{Error: fmt.Errorf("decode line failure %w", err)}
			return
		}
		s.stream <- types.PortStream{ID: fmt.Sprintf("%s", tkn), Port: port}
	}

	if _, err := decoder.Token(); err != nil {
		s.stream <- types.PortStream{Error: fmt.Errorf("decode closing delimiter: %w", err)}
		return
	}
}
