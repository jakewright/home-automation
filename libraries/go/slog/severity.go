package slog

// Severity is a subset of the syslog severity levels
type Severity string

const (
	// DebugSeverity is the severity used for debug-level messages
	DebugSeverity Severity = "DEBUG"

	// InfoSeverity is the severity used for informational messages
	InfoSeverity Severity = "INFO"

	// WarnSeverity is the severity used for warning conditions
	WarnSeverity Severity = "WARNING"

	// ErrorSeverity is the severity used for error conditions
	ErrorSeverity Severity = "ERR"
)
