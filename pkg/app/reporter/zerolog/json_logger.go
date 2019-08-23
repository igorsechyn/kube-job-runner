package zerolog

import (
	"os"

	"github.com/rs/zerolog"
)

// NewJSONLogger returns a new instance
func NewJSONLogger() *JSONLogger {
	log := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return &JSONLogger{
		l: log,
	}
}

// JSONLogger prints messages in json format
type JSONLogger struct {
	l zerolog.Logger
}

// Info channel for output
func (logger *JSONLogger) Info(message string, fields map[string]interface{}) {
	log := logger.withFields(fields)
	log.Info().Msg(message)
}

// Error channel for output
func (logger *JSONLogger) Error(err error, fields map[string]interface{}) {
	log := logger.withFields(fields)

	log.Error().Err(err).Msg(err.Error())
}

func (logger *JSONLogger) withFields(fields map[string]interface{}) zerolog.Logger {
	log := logger.l.With().Logger()
	for key, value := range fields {
		log = log.With().Interface(key, value).Logger()
	}

	return log
}
