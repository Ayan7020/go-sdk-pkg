package logger

import (
	"io"
	"sync/atomic"

	otellog "go.opentelemetry.io/otel/log"
)

type Fields struct {
	Key   string
	Value any
}

type Logger interface {
	Info()
	// With(...Fields) Logger
	// WithOtel()
}

type Settings struct {
	Service, Version, Env string 
	// Level                 any
	OtelLogs otellog.LoggerProvider
	Output   io.Writer
}
type Factory func(Settings) (Logger, error)

var def atomic.Pointer[Logger]

func SetDefault(l Logger) { def.Store(&l) }
func Default() Logger {
	logger := def.Load()
	return *logger
}
