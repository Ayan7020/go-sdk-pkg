package observability

import (
	"context"
	"errors"
	"os"
	"sync"

	"github.com/Ayan7020/go-sdk-pkg/observability/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

type Shutdown func(context.Context) error

func Init(ctx context.Context, cfg Config, opts ...Option) (Shutdown, error) {
	if cfg.ServiceName == "" {
		return nil, errors.New("observability: ServiceName is required")
	}

	s := defaultSettings(cfg)
	for _, opt := range opts {
		err := opt.apply(&s)
		if err != nil {
			return nil, err
		}
	}

	res, err := resolveResouce(ctx, s)
	if err != nil {
		return nil, err
	}

	lp, err := resolveLogger(ctx, s, res)
	if err != nil {
		return nil, err
	}

	log, err := s.logger(logger.Settings{
		Service:  s.ServiceName,
		Version:  s.Version,
		Env:      s.Environment,
		Output:   os.Stdout,
		OtelLogs: lp,
	})
	if err != nil {
		return nil, err
	}
	logger.SetDefault(log)

	traceProvider, err := resolveTracing(ctx, s, res)
	if err != nil {
		return nil, err
	}

	meterProvider, err := resolveMetrics(ctx, s, res)
	if err != nil {
		_ = traceProvider.Shutdown(ctx)
		if provider, ok := lp.(interface{ Shutdown(context.Context) error }); ok {
			_ = provider.Shutdown(ctx)
		}
		return nil, err
	}
	otel.SetMeterProvider(meterProvider)

	var shutdownOnce sync.Once
	var shutdownErr error
	return func(shutdownCtx context.Context) error {
		shutdownOnce.Do(func() {
			if err := traceProvider.Shutdown(shutdownCtx); err != nil {
				shutdownErr = errors.Join(shutdownErr, err)
			}
			if err := meterProvider.Shutdown(shutdownCtx); err != nil {
				shutdownErr = errors.Join(shutdownErr, err)
			}
			if provider, ok := lp.(interface{ Shutdown(context.Context) error }); ok {
				if err := provider.Shutdown(shutdownCtx); err != nil {
					shutdownErr = errors.Join(shutdownErr, err)
				}
			}
		})
		return shutdownErr
	}, nil
}

func Meter(name string) metric.Meter {
	return otel.Meter(name)
}
