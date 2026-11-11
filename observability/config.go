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
	ExporterEndpoint string
	Insecure         bool
	logger           logger.Factory
}

func defaultSettings(c Config) settings {
	return settings{
		Config:           c,
		logger:           logger.NewZapJSON,
		Insecure:         false,
		ExporterEndpoint: os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"),
	}
}
