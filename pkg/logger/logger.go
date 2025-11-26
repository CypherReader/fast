package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Init initializes the global logger with the specified log level
func Init(level string) {
	// Configure zerolog
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Set log level
	switch level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// Pretty print in development
	if os.Getenv("GIN_MODE") != "release" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

// Info returns a new info level event
func Info() *zerolog.Event {
	return log.Info()
}

// Debug returns a new debug level event
func Debug() *zerolog.Event {
	return log.Debug()
}

// Warn returns a new warn level event
func Warn() *zerolog.Event {
	return log.Warn()
}

// Error returns a new error level event
func Error() *zerolog.Event {
	return log.Error()
}

// Fatal returns a new fatal level event (exits after logging)
func Fatal() *zerolog.Event {
	return log.Fatal()
}
