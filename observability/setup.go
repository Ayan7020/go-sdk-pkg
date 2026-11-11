package observability

import (
	"context"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	otellog "go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
)

func resolveResouce(ctx context.Context, s settings) (*resource.Resource, error) {
	r := resource.Default()

	CtxRes, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceName(s.ServiceName),
			semconv.ServiceVersion(s.Version),
			semconv.DeploymentEnvironmentName(s.Environment),
			semconv.ServiceNamespace(s.ServiceNamespace),
			semconv.ServiceInstanceID(s.ServiceInstanceID),
			semconv.CloudRegion(s.DeploymentRegion),
		),
		resource.WithFromEnv(),
		resource.WithProcess(),
	)

	r, err = resource.Merge(r, CtxRes)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func resolveTracing(ctx context.Context, s settings, res *resource.Resource) (*sdktrace.TracerProvider, error) {
	if res == nil {
		return nil, fmt.Errorf("failed to resolve tracing: resources are not initialized")
	}

	var exporter sdktrace.SpanExporter
	var err error
	if strings.HasPrefix(s.ExporterEndpoint, "http") {
		options := []otlptracehttp.Option{otlptracehttp.WithEndpointURL(s.ExporterEndpoint)}
		if s.Insecure {
			options = append(options, otlptracehttp.WithInsecure())
		}
		exporter, err = otlptracehttp.New(ctx, options...)
	} else {
		options := []otlptracegrpc.Option{otlptracegrpc.WithEndpoint(s.ExporterEndpoint)}
		if s.Insecure {
			options = append(options, otlptracegrpc.WithInsecure())
		}
		exporter, err = otlptracegrpc.New(ctx, options...)
	}
	if err != nil {
		return nil, err
	}

	batchOptions := make([]sdktrace.BatchSpanProcessorOption, 0, 2)
	if s.batchTimeout > 0 {
		batchOptions = append(batchOptions, sdktrace.WithBatchTimeout(s.batchTimeout))
	}
	if s.exportTimeout > 0 {
		batchOptions = append(batchOptions, sdktrace.WithExportTimeout(s.exportTimeout))
	}

	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter, batchOptions...),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(s.sampler),
	)

	otel.SetTracerProvider(traceProvider)
	return traceProvider, nil
}

func resolveLogger(ctx context.Context, s settings, res *resource.Resource) (otellog.LoggerProvider, error) {
	if s.ExporterEndpoint == "" {
		return nil, nil
	}
	var (
		exp sdklog.Exporter
		err error
	)

	if strings.HasPrefix(s.ExporterEndpoint, "http") {
		options := []otlploghttp.Option{otlploghttp.WithEndpointURL(s.ExporterEndpoint)}
		if s.Insecure {
			options = append(options, otlploghttp.WithInsecure())
		}
		exp, err = otlploghttp.New(ctx, options...)
	} else {
		options := []otlploggrpc.Option{otlploggrpc.WithEndpoint(s.ExporterEndpoint)}
		if s.Insecure {
			options = append(options, otlploggrpc.WithInsecure())
		}
		exp, err = otlploggrpc.New(ctx, options...)
	}

	if err != nil {
		return nil, fmt.Errorf("otlp log exporter: %w", err)
	}

	lp := sdklog.NewLoggerProvider(
		sdklog.WithResource(res),
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exp)),
	)
	return lp, nil
}

func resolveMetrics(ctx context.Context, s settings, res *resource.Resource) (*sdkmetric.MeterProvider, error) {
	if res == nil {
		return nil, fmt.Errorf("observability: resources are not initialized")
	}

	if s.ExporterEndpoint == "" {
		return sdkmetric.NewMeterProvider(sdkmetric.WithResource(res)), nil
	}

	var (
		exporter sdkmetric.Exporter
		err      error
	)
	if strings.HasPrefix(s.ExporterEndpoint, "http") {
		options := []otlpmetrichttp.Option{otlpmetrichttp.WithEndpointURL(s.ExporterEndpoint)}
		if s.Insecure {
			options = append(options, otlpmetrichttp.WithInsecure())
		}
		exporter, err = otlpmetrichttp.New(ctx, options...)
	} else {
		options := []otlpmetricgrpc.Option{otlpmetricgrpc.WithEndpoint(s.ExporterEndpoint)}
		if s.Insecure {
			options = append(options, otlpmetricgrpc.WithInsecure())
		}
		exporter, err = otlpmetricgrpc.New(ctx, options...)
	}
	if err != nil {
		return nil, fmt.Errorf("otlp metric exporter: %w", err)
	}

	readerOptions := make([]sdkmetric.PeriodicReaderOption, 0, 2)
	if s.metricInterval > 0 {
		readerOptions = append(readerOptions, sdkmetric.WithInterval(s.metricInterval))
	}
	if s.exportTimeout > 0 {
		readerOptions = append(readerOptions, sdkmetric.WithTimeout(s.exportTimeout))
	}

	return sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter, readerOptions...)),
	), nil
}
