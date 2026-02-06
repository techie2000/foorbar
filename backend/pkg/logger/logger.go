package logger

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var logger zerolog.Logger

// Init initializes the logger
func Init(level string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	logLevel := zerolog.InfoLevel
	switch level {
	case "debug":
		logLevel = zerolog.DebugLevel
	case "warn":
		logLevel = zerolog.WarnLevel
	case "error":
		logLevel = zerolog.ErrorLevel
	}

	logger = zerolog.New(os.Stdout).
		Level(logLevel).
		With().
		Timestamp().
		Caller().
		Logger()

	log.Logger = logger
}

// Debug returns a debug level event
func Debug() *zerolog.Event {
	return logger.Debug()
}

// Info returns an info level event
func Info() *zerolog.Event {
	return logger.Info()
}

// Warn returns a warn level event
func Warn() *zerolog.Event {
	return logger.Warn()
}

// Error returns an error level event
func Error() *zerolog.Event {
	return logger.Error()
}

// Fatal returns a fatal level event
func Fatal() *zerolog.Event {
	return logger.Fatal()
}
