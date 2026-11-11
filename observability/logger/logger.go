package logger

import (
	"context"
	"io"
	"os"
	"sync"
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
var fallbackOnce sync.Once

func SetDefault(l Logger) { def.Store(&l) }

func Default() Logger {
	if logger := def.Load(); logger != nil && *logger != nil {
		return *logger
	}

	fallbackOnce.Do(func() {
		fallback, err := NewZapConsole(Settings{Output: os.Stderr})
		if err == nil {
			def.Store(&fallback)
		}
	})
	if logger := def.Load(); logger != nil && *logger != nil {
		return *logger
	}
	return noopLogger{}
}

func FromContext(ctx context.Context) Logger {
	log := Default()

	if ctx == nil {
		return log
	}
	if spanContext := trace.SpanContextFromContext(ctx); spanContext.IsValid() {
		log = log.With(
			Fields{Key: "trace_id", Value: spanContext.TraceID().String()},
			Fields{Key: "span_id", Value: spanContext.SpanID().String()},
		)
	}

	return log
}

type noopLogger struct{}

func (noopLogger) Debug(string, ...Fields) {}
func (noopLogger) Info(string, ...Fields)  {}
func (noopLogger) Warn(string, ...Fields)  {}
func (noopLogger) Error(string, ...Fields) {}
func (n noopLogger) With(...Fields) Logger { return n }
