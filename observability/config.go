package observability

import (
	"os"
	"time"

	"github.com/Ayan7020/go-sdk-pkg/observability/logger"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type Config struct {
	ServiceName, Version, Environment                     string
	ServiceNamespace, ServiceInstanceID, DeploymentRegion string
	ExporterEndpoint                                      string
	Insecure                                              bool
}

type settings struct {
	Config 
	logger           logger.Factory
	sampler          sdktrace.Sampler
	batchTimeout     time.Duration
	exportTimeout    time.Duration
	metricInterval   time.Duration
}

func defaultSettings(c Config) settings {
	if c.ExporterEndpoint == "" {
		c.ExporterEndpoint = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	}

	return settings{
		Config:           c,
		logger:           logger.NewZapJSON,  
		sampler:          sdktrace.ParentBased(sdktrace.TraceIDRatioBased(1)),
	}
}
