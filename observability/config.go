package observability

import (
	"os"

	"github.com/Ayan7020/go-sdk-pkg/observability/logger"
)

type Config struct {
	ServiceName, Version, Environment string
}

type settings struct {
	Config
	logger   logger.Factory
	// logLevel any
	exporter bool
}

func defaultSettings(c Config) settings {
	return settings{
		Config:   c,
		logger:   logger.NewZapJSON, // n 
		// logLevel: logger.LevelInfo,
		exporter: os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") != "",
	}
}
