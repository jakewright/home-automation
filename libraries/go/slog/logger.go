package slog

import "os"

// Logger logs logs
type Logger interface {
	//Log(severity Severity, format string, params ...interface{})
	Log(*Log)
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
func Debug(format string, params ...interface{}) {
	mustGetDefaultLogger().Log(newFromFormat(DebugSeverity, format, params...))
}

// Info logs with INFO severity
func Info(format string, params ...interface{}) {
	mustGetDefaultLogger().Log(newFromFormat(InfoSeverity, format, params...))
}

// Warn logs with WARNING severity
func Warn(format string, params ...interface{}) {
	mustGetDefaultLogger().Log(newFromFormat(WarnSeverity, format, params...))
}

// Error logs with ERR severity
func Error(format string, params ...interface{}) {
	mustGetDefaultLogger().Log(newFromFormat(ErrorSeverity, format, params...))
}

// Fatal logs with ERR severity and terminates the program
func Fatal(format string, params ...interface{}) {
	Error(format, params...)
	os.Exit(1)
}
