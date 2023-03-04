package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

type FatalLogger interface {
	Logger
	Fatalf(format string, a ...any)
}

type Logger interface {
	Debugf(format string, a ...any)
	Infof(format string, a ...any)
	Warnf(format string, a ...any)
	Errorf(format string, a ...any)
}

func NewLogger(opts ...Option) FatalLogger {
	l := logger{
		writers: map[Severity]io.Writer{
			Debug:   os.Stdout,
			Info:    os.Stdout,
			Warning: os.Stdout,
			Error:   os.Stderr,
			Fatal:   os.Stderr,
		},
	}
	for _, opt := range opts {
		opt(&l)
	}
	return &l
}

type logger struct {
	writers      map[Severity]io.Writer
	outputFormat uint8
}

var _ = (*logger)(nil)

const (
	Plain uint8 = iota
	JSON
)

func (l *logger) Debugf(format string, a ...any) {
	l.print(Debug, format, a...)
}

func (l *logger) Infof(format string, a ...any) {
	l.print(Info, format, a...)
}

func (l *logger) Warnf(format string, a ...any) {
	l.print(Warning, format, a...)
}

func (l *logger) Errorf(format string, a ...any) {
	l.print(Error, format, a...)
}

func (l *logger) Fatalf(format string, a ...any) {
	l.print(Fatal, format, a...)
	os.Exit(1)
}

func (l *logger) print(severity Severity, format string, a ...any) {
	msg := message{
		Time:     time.Now(),
		Severity: severity,
		Message:  fmt.Sprintf(format, a...),
	}

	var err error
	switch l.outputFormat {
	case Plain:
		_, err = l.writers[severity].Write([]byte(fmt.Sprintf(
			"%s %s %s", msg.Time.Format(time.RFC3339), msg.Severity, msg.Message,
		)))
	case JSON:
		err = json.NewEncoder(l.writers[severity]).Encode(msg)
	}

	if err != nil {
		// TODO: do anything?
	}
}
