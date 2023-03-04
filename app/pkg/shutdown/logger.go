package shutdown

import "log"

type Logger interface {
	Infof(format string, a ...any)
	Errorf(format string, a ...any)
}

var defaultLogger Logger = &logger{}

type logger struct{}

var _ = (*logger)(nil)

func (l *logger) Infof(format string, a ...any) {
	log.Printf("INFO: "+format+"\n", a...)
}

func (l *logger) Errorf(format string, a ...any) {
	log.Printf("ERROR: "+format+"\n", a...)
}
