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

// Debug logs with DEBUG severity
func Debug(format string, a ...interface{}) {
	mustGetDefaultLogger().Log(newEventFromFormat(DebugSeverity, format, a...))
}

// Info logs with INFO severity
func Info(format string, a ...interface{}) {
	mustGetDefaultLogger().Log(newEventFromFormat(InfoSeverity, format, a...))
}

// Warn logs with WARNING severity
func Warn(format string, a ...interface{}) {
	mustGetDefaultLogger().Log(newEventFromFormat(WarnSeverity, format, a...))
}

// Error logs with ERROR severity
func Error(format string, a ...interface{}) {
	mustGetDefaultLogger().Log(newEventFromFormat(ErrorSeverity, format, a...))
}

// Panic logs with ERROR severity and then panics
func Panic(format string, a ...interface{}) {
	Error(format, a...)
	panic(newEventFromFormat(ErrorSeverity, format, a...))
}
