package shutdown

import "time"

type Option func(config *settings)

func WithTimeout(timeout time.Duration) Option { return func(s *settings) { s.timeout = timeout } }

func WithLogger(logger Logger) Option { return func(s *settings) { s.logger = logger } }

type settings struct {
	timeout time.Duration
	logger  Logger
}
