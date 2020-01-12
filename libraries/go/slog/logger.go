package slog

// Logger logs logs
type Logger interface {
	Log(*Event)
}

// DefaultLogger should be used to log all events
var DefaultLogger Logger

func mustGetDefaultLogger() Logger {
	if DefaultLogger == nil {
		DefaultLogger = NewStdoutLogger()
	}

	return DefaultLogger
}

// Debugf logs with DEBUG severity
func Debugf(format string, a ...interface{}) {
	mustGetDefaultLogger().Log(newEventFromFormat(DebugSeverity, format, a...))
}

// Infof logs with INFO severity
func Infof(format string, a ...interface{}) {
	mustGetDefaultLogger().Log(newEventFromFormat(InfoSeverity, format, a...))
}

// Warnf logs with WARNING severity
func Warnf(format string, a ...interface{}) {
	mustGetDefaultLogger().Log(newEventFromFormat(WarnSeverity, format, a...))
}

// Errorf logs with ERROR severity
func Errorf(format string, a ...interface{}) {
	mustGetDefaultLogger().Log(newEventFromFormat(ErrorSeverity, format, a...))
}

// Error logs with ERROR severity
func Error(v interface{}) {
	mustGetDefaultLogger().Log(newEvent(ErrorSeverity, v))
}

// Panicf logs with ERROR severity and then panics
func Panicf(format string, a ...interface{}) {
	Errorf(format, a...)
	panic(newEventFromFormat(ErrorSeverity, format, a...))
}

// Panic logs with ERROR severity and then panics
func Panic(v interface{}) {
	Error(v)
	panic(newEvent(ErrorSeverity, v))
}
