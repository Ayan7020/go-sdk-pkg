package observability

import "github.com/Ayan7020/go-sdk-pkg/observability/logger"

type Option func(*settings) error

func (opt Option) apply(s *settings) error {
	return opt(s)
}

func WithZapJSONLogger() Option {
	return func(s *settings) error { s.logger = logger.NewZapJSON; return nil }
}
func WithZapConsoleLogger() Option {
	return func(s *settings) error { s.logger = logger.NewZapConsole; return nil }
}

func WithExporterEndpoint(endpoint string) Option {
	return func(s *settings) error { s.ExporterEndpoint = endpoint; return nil }
}

func WithInSecure() Option {
	return func(s *settings) error { s.Insecure = true; return nil }
} 