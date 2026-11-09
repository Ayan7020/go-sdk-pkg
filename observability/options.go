package observability

import "github.com/kalento/go_api/pkg/observability/logger"

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
 

// func WithLogLevel(l logger.Level) Option {
// 	return func(s *settings) error { s.logLevel = l; return nil }
// }
// func WithoutExporter() Option { // tests, local debugging
// 	return func(s *settings) error { s.exporter = false; return nil }
// }
