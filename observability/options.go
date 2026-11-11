package observability

import (
	"fmt"
	"time"

	"github.com/Ayan7020/go-sdk-pkg/observability/logger"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Option func(*settings) error

func (opt Option) apply(s *settings) error {
	return opt(s)
}

func WithZapJSONLogger() Option {
	return func(s *settings) error { s.logger = logger.NewZapJSON; return nil }
}
func WithZapConsoleLogger() Option {
	return func(s *settings) error { s.logger = logger.NewZapConsole; return nil }
}

func WithExporterEndpoint(endpoint string) Option {
	return func(s *settings) error {
		if endpoint == "" {
			return fmt.Errorf("observability: exporter endpoint cannot be empty")
		}
		s.ExporterEndpoint = endpoint
		return nil
	}
}

func WithInsecure() Option {
	return func(s *settings) error { s.Insecure = true; return nil }
}

func WithInSecure() Option { return WithInsecure() }

func WithTraceSampleRatio(ratio float64) Option {
	return func(s *settings) error {
		if ratio < 0 || ratio > 1 {
			return fmt.Errorf("observability: trace sample ratio must be between 0 and 1")
		}
		s.sampler = sdktrace.ParentBased(sdktrace.TraceIDRatioBased(ratio))
		return nil
	}
}

func WithBatchTimeout(timeout time.Duration) Option {
	return func(s *settings) error {
		if timeout <= 0 {
			return fmt.Errorf("observability: batch timeout must be positive")
		}
		s.batchTimeout = timeout
		return nil
	}
}

func WithExportTimeout(timeout time.Duration) Option {
	return func(s *settings) error {
		if timeout <= 0 {
			return fmt.Errorf("observability: export timeout must be positive")
		}
		s.exportTimeout = timeout
		return nil
	}
}

func WithMetricInterval(interval time.Duration) Option {
	return func(s *settings) error {
		if interval <= 0 {
			return fmt.Errorf("observability: metric interval must be positive")
		}
		s.metricInterval = interval
		return nil
	}
}
