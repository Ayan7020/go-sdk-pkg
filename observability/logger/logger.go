package logger

import (
	"context"
	"io"
	"sync/atomic"

	otellog "go.opentelemetry.io/otel/log"
	"go.opentelemetry.io/otel/trace"
)

type Fields struct {
	Key   string
	Value any
}

type Logger interface {
	Debug(msg string, field ...Fields)
	Info(msg string, field ...Fields)
	Warn(msg string, field ...Fields)
	Error(msg string, field ...Fields)
	With(...Fields) Logger 
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

// Todo: add fallback logger
func Default() Logger {
	logger := def.Load()
	return *logger
}

func FromContext(ctx context.Context) Logger {
	log := Default()

	if spanContext := trace.SpanContextFromContext(ctx); spanContext.IsValid() {
		log = log.With(
			Fields{Key: "trace_id", Value: spanContext.TraceID().String()},
			Fields{Key: "span_id", Value: spanContext.SpanID().String()},
		)
	}

	return log
}
