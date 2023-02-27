package logger

var defaultLogger Logger = &logger{outputFormat: JSON}

func Debugf(format string, a ...any) { defaultLogger.Debugf(format, a...) }
func Infof(format string, a ...any)  { defaultLogger.Infof(format, a...) }
func Warnf(format string, a ...any)  { defaultLogger.Warnf(format, a...) }
func Errorf(format string, a ...any) { defaultLogger.Errorf(format, a...) }
