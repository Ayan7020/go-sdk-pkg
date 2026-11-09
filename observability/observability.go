package observability

import (
	"context"
	"errors"
	"os"

	"github.com/kalento/go_api/pkg/observability/logger"
)

func Init(ctx context.Context, cfg Config, opts ...Option) error {
	if cfg.ServiceName == "" {
		return errors.New("observability: ServiceName is required")
	}

	s := defaultSettings(cfg)
	for _, opt := range opts {
		err := opt.apply(&s)
		if err != nil {
			return err
		}
	}
  
	res, err := resolveResouce(ctx, s)
	if err != nil {
		return err
	}

	lp,err := resolveLogger(ctx,s,res)
	if err != nil {
		return err
	}

	log, err := s.logger(logger.Settings{
		Service: s.ServiceName,
		Version: s.Version,
		Env:     s.Environment,
		Output:  os.Stdout,
		OtelLogs: lp,
	})
	if err != nil {
		return err
	}
	logger.SetDefault(log)
	
	if err := resolveTracing(ctx, s, res); err != nil {
		return err
	}

	return nil
}
