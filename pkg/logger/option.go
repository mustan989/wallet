package logger

import "io"

type Option func(logger *logger)

func WithSeverityWriter(severity Severity, writer io.Writer) Option {
	return func(logger *logger) { logger.writers[severity] = writer }
}

func WithWriters(writers map[Severity]io.Writer) Option {
	return func(logger *logger) {
		for severity, writer := range writers {
			logger.writers[severity] = writer
		}
	}
}

func WithOutputFormat(format uint8) Option {
	return func(logger *logger) { logger.outputFormat = format }
}
