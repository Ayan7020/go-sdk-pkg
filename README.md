# Go Observability SDK

Reusable observability infrastructure for Go services, built on
[OpenTelemetry](https://opentelemetry.io/) and
[Zap](https://github.com/uber-go/zap).

The SDK centralizes telemetry setup so applications can configure logs, traces,
and metrics consistently without duplicating exporter and resource setup in
every service.

## Features

- OpenTelemetry traces, metrics, and logs
- OTLP HTTP and gRPC exporter support
- Centralized endpoint and insecure transport configuration
- Structured JSON and console logging with Zap
- Trace and span IDs automatically attached to context-aware logs
- Configurable trace sampling
- Configurable trace batching and metric export intervals
- Standard service metadata: version, environment, namespace, instance, and region
- Graceful and idempotent shutdown for flushing telemetry
- Safe logger fallback before explicit initialization

## Installation

```bash
go get github.com/Ayan7020/go-sdk-pkg
```

## Quick start

```go
package main

import (
	"context"
	"log"

	"github.com/Ayan7020/go-sdk-pkg/observability"
)

func main() {
	ctx := context.Background()

	shutdown, err := observability.Init(ctx, observability.Config{
		ServiceName:      "orders-api",
		Version:          "1.0.0",
		Environment:      "production",
		ExporterEndpoint: "otel-collector:4317",
	}, observability.WithTraceSampleRatio(0.25))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := shutdown(context.Background()); err != nil {
			log.Printf("observability shutdown failed: %v", err)
		}
	}()
}
```

`ServiceName` is required. If `ExporterEndpoint` is omitted, the SDK reads
`OTEL_EXPORTER_OTLP_ENDPOINT`. If neither is configured, the SDK still
initializes local tracer and meter providers without exporting telemetry.

Endpoints beginning with `http://` or `https://` use OTLP/HTTP; host-and-port
endpoints use OTLP/gRPC. Use `Insecure: true` or `WithInsecure()` only for
plaintext collector connections.

For a plaintext collector connection, set `Insecure: true`:

```go
observability.Config{
	ServiceName:      "orders-api",
	ExporterEndpoint: "localhost:4317",
	Insecure:         true,
}
```

## Logging

The SDK configures a global structured logger. Use `FromContext` when handling
a request or operation that already has an OpenTelemetry span:

```go
import "github.com/Ayan7020/go-sdk-pkg/observability/logger"

log := logger.FromContext(ctx)
log.Info("order created",
	logger.Fields{Key: "order_id", Value: orderID},
)
```

`FromContext` returns a child logger containing valid `trace_id` and `span_id`
fields without mutating the global logger.

Available log levels:

```go
log.Debug("debug details")
log.Info("operation started")
log.Warn("operation is slow")
log.Error("operation failed", logger.Fields{Key: "error", Value: err})
```

Choose console logging for local development:

```go
shutdown, err := observability.Init(
	ctx,
	cfg,
	observability.WithZapConsoleLogger(),
)
```

## Metrics

The SDK configures the global OpenTelemetry meter provider. Applications define
their own domain metrics:

```go
meter := observability.Meter("orders")

created, err := meter.Int64Counter("orders.created")
if err != nil {
	log.Fatal(err)
}

created.Add(ctx, 1)
```

## Configuration options

```go
shutdown, err := observability.Init(
	ctx,
	cfg,
	observability.WithTraceSampleRatio(0.25),
	observability.WithBatchTimeout(5*time.Second),
	observability.WithExportTimeout(10*time.Second),
	observability.WithMetricInterval(15*time.Second),
)
```

Available options include:

- `WithExporterEndpoint`
- `WithInsecure`
- `WithZapJSONLogger`
- `WithZapConsoleLogger`
- `WithTraceSampleRatio`
- `WithBatchTimeout`
- `WithExportTimeout`
- `WithMetricInterval`

## Design principles

- Keep application metrics in the application, not in the SDK.
- Keep the core package independent of web frameworks.
- Make telemetry configuration explicit and reusable.
- Preserve request safety by deriving child loggers instead of mutating globals.
- Flush telemetry during graceful application shutdown.

## Project status

This project is actively evolving toward a production-ready observability
foundation for Go applications. New improvements, tests, and integrations are
added regularly.

## Validation

```bash
go test ./...
go vet ./...
```

## License

License information will be added in a future release.
