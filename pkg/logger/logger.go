package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/shanto-323/chat-ai/config"
)

func New(cfg *config.Config) *zerolog.Logger {
	var logLevel zerolog.Level
	level := cfg.Logging.Level

	switch level {
	case "debug":
		logLevel = zerolog.DebugLevel
	case "info":
		logLevel = zerolog.InfoLevel
	case "warn":
		logLevel = zerolog.WarnLevel
	case "error":
		logLevel = zerolog.ErrorLevel
	default:
		logLevel = zerolog.InfoLevel
	}

	zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	var writer io.Writer
	baseWriter := os.Stdout

	if cfg.IsProd() {
		err := os.MkdirAll("/var/lib/logs", 0o755)
		if err == nil {
			file, err := os.OpenFile("/var/lib/logs/app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
			if err == nil {
				writer = io.Writer(file)
			} else {
				file.Close()
				writer = io.Writer(baseWriter)
			}
		} else {
			writer = io.Writer(baseWriter)
		}
	} else {
		writer = io.Writer(baseWriter)
	}

	logger := zerolog.New(writer).
		Level(logLevel).
		With().
		Timestamp().
		Str("service", cfg.Primary.ServiceName).
		Str("environment", cfg.Primary.Env).
		Logger()

	return &logger
}
