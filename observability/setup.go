package observability

import (
	"context"
	"fmt"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	otellog "go.opentelemetry.io/otel/log"
	sdklog "go.opentelemetry.io/otel/sdk/log"
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

func resolveTracing(ctx context.Context, s settings, res *resource.Resource) error {
	if res == nil {
		return fmt.Errorf("failed to resolve Tracing! first initialized the resources")
	}

	var exporter sdktrace.SpanExporter
	var err error
	options := make([]otlptracegrpc.Option, 0, 2)
	if strings.HasPrefix(s.ExporterEndpoint, "http") {
		options = append(options, otlptracegrpc.WithEndpointURL(s.ExporterEndpoint))
	} else if s.ExporterEndpoint != "" {
		options = append(options, otlptracegrpc.WithEndpoint(s.ExporterEndpoint))
	}
	if s.Insecure {
		options = append(options, otlptracegrpc.WithInsecure())
	}
	if len(options) > 0 {
		exporter, err = otlptracegrpc.New(ctx, options...)
	} else {
		exporter, err = otlptracegrpc.New(ctx)
	}
	if err != nil {
		return err
	}

	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)

	otel.SetTracerProvider(traceProvider)
	return nil
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
		options := make([]otlploghttp.Option, 0, 2)
		options = append(options, otlploghttp.WithEndpointURL(s.ExporterEndpoint))
		if s.Insecure {
			options = append(options, otlploghttp.WithInsecure())
		}
		exp, err = otlploghttp.New(ctx, options...)
	} else if s.ExporterEndpoint != "" {
		options := make([]otlploggrpc.Option, 0, 2)
		options = append(options, otlploggrpc.WithEndpoint(s.ExporterEndpoint))
		if s.Insecure {
			options = append(options, otlploggrpc.WithInsecure())
		}
		exp, err = otlploggrpc.New(ctx, options...)
	} else {
		exp, err = otlploggrpc.New(ctx)
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
