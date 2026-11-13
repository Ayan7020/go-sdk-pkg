package observability

import (
	"context"
	"testing"
)

func TestInitWithoutExporter(t *testing.T) {
	shutdown, err := Init(context.Background(), Config{ServiceName: "test-service"})
	if err != nil {
		t.Fatalf("Init() error = %v", err)
	}

	if shutdown == nil {
		t.Fatal("Init() returned a nil shutdown function")
	}
	if err := shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown() error = %v", err)
	}
	if err := shutdown(context.Background()); err != nil {
		t.Fatalf("second shutdown() error = %v", err)
	}
}

func TestInitRejectsInvalidInput(t *testing.T) {
	tests := []struct {
		name string
		cfg  Config
		opts []Option
	}{
		{
			name: "missing service name",
		},
		{
			name: "nil option",
			cfg:  Config{ServiceName: "test-service"},
			opts: []Option{nil},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := Init(context.Background(), tt.cfg, tt.opts...); err == nil {
				t.Fatal("Init() error = nil, want error")
			}
		})
	}
}

func TestWithExporterEndpointTrimsAndValidates(t *testing.T) {
	settings := defaultSettings(Config{})
	if err := WithExporterEndpoint("  localhost:4317  ").apply(&settings); err != nil {
		t.Fatalf("WithExporterEndpoint() error = %v", err)
	}
	if settings.ExporterEndpoint != "localhost:4317" {
		t.Fatalf("ExporterEndpoint = %q, want %q", settings.ExporterEndpoint, "localhost:4317")
	}

	if err := WithExporterEndpoint("  ").apply(&settings); err == nil {
		t.Fatal("WithExporterEndpoint() error = nil, want error")
	}
}

func TestIsHTTPEndpoint(t *testing.T) {
	for _, endpoint := range []string{"http://localhost:4318", "HTTPS://collector.example"} {
		if !isHTTPEndpoint(endpoint) {
			t.Errorf("isHTTPEndpoint(%q) = false, want true", endpoint)
		}
	}
	if isHTTPEndpoint("localhost:4317") {
		t.Error("isHTTPEndpoint(gRPC endpoint) = true, want false")
	}
}
