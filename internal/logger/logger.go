package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewZapLogger initialise zap logger
func NewZapLogger(appName string, level zapcore.Level) *zap.Logger {
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "ts_app"
	config.EncoderConfig.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	config.Level.SetLevel(level)

	logger, err := config.Build()
	if err != nil {
		log.Fatalf("initialise zap logger: %v", err)
	}
	return logger.With(zap.String("app", appName))
}
