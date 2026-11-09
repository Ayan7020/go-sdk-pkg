package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"go.opentelemetry.io/contrib/bridges/otelzap"
)

type zapLogger struct {
	l *zap.Logger
}

func NewZapJSON(s Settings) (Logger, error)    { return newZap(s, false) }
func NewZapConsole(s Settings) (Logger, error) { return newZap(s, true) }

func newZap(s Settings, console bool) (Logger, error) {
	out := s.Output
	if out == nil {
		out = os.Stdout
	}

	var enc zapcore.Encoder
	if console {
		ec := zap.NewDevelopmentEncoderConfig()
		ec.EncodeLevel = zapcore.CapitalColorLevelEncoder
		ec.EncodeTime = zapcore.ISO8601TimeEncoder
		enc = zapcore.NewConsoleEncoder(ec)
	} else {
		ec := zap.NewProductionEncoderConfig()
		ec.TimeKey = "time"
		ec.EncodeTime = zapcore.ISO8601TimeEncoder
		enc = zapcore.NewJSONEncoder(ec)
	}

	core := zapcore.NewCore(enc, zapcore.AddSync(out), zapcore.DebugLevel)

	var base []zap.Field
	if s.Service != "" {
		base = append(base, zap.String("service", s.Service))
	}
	if s.Version != "" {
		base = append(base, zap.String("version", s.Version))
	}
	if s.Env != "" {
		base = append(base, zap.String("env", s.Env))
	}

	if s.OtelLogs != nil {
		core = zapcore.NewTee(core, otelzap.NewCore(
			s.Service,
			otelzap.WithVersion(s.Version),
			otelzap.WithLoggerProvider(s.OtelLogs),
		))
	}

	return &zapLogger{l: zap.New(
		core,
		zap.AddCaller(),
		zap.AddCallerSkip(2),
		zap.AddStacktrace(zapcore.ErrorLevel))}, nil
}

func (z *zapLogger) Info(msg string, field ...Fields) { z.l.Info(msg, mapToZapFields(field...)...) }

func mapToZapFields(field ...Fields) []zap.Field {
	if len(field) == 0 {
		return nil
	}

	out := make([]zap.Field, len(field))
	for i, f := range field {
		switch v := f.Value.(type) {
		case string:
			out[i] = zap.String(f.Key, v)
		case int:
			out[i] = zap.Int(f.Key, v)
		case bool:
			out[i] = zap.Bool(f.Key, v)
		}
	}

	return out
}
