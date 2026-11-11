# Go Observability SDK

This package provides centralized OpenTelemetry traces, metrics, and logs with
structured Zap logging.

## Initialize

```go
shutdown, err := observability.Init(ctx, observability.Config{
	ServiceName:      "orders-api",
	Version:          "1.0.0",
	Environment:      "production",
	ExporterEndpoint: "otel-collector:4317",
}, observability.WithTraceSampleRatio(0.25))
if err != nil {
	return err
}
defer shutdown(context.Background())
```

`ExporterEndpoint` may also be supplied through
`OTEL_EXPORTER_OTLP_ENDPOINT`. Set `Insecure: true` for a plaintext collector
connection.

## Metrics

The SDK configures the global meter provider. Applications define their own
business metrics:

```go
meter := observability.Meter("orders")
created, err := meter.Int64Counter("orders.created")
if err != nil {
	return err
}
created.Add(ctx, 1)
```

## Logging

```go
log := logger.FromContext(ctx)
log.Info("order created",
	logger.Fields{Key: "order_id", Value: orderID},
)
```

`FromContext` adds valid OpenTelemetry `trace_id` and `span_id` values without
mutating the global logger.
