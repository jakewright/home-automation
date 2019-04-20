package slog

// Severity is a subset of the syslog severity levels
type Severity int

const (
	// DebugSeverity is the severity used for debug-level messages
	DebugSeverity Severity = 2

	// InfoSeverity is the severity used for informational messages
	InfoSeverity Severity = 3

	// WarnSeverity is the severity used for warning conditions
	WarnSeverity Severity = 5

	// ErrorSeverity is the severity used for error conditions
	ErrorSeverity Severity = 6
)

func (s Severity) String() string {
	switch s {
	case DebugSeverity:
		return "DEBUG"
	case InfoSeverity:
		return "INFO"
	case WarnSeverity:
		return "WARN"
	case ErrorSeverity:
		return "ERROR"
	}

	return "UNKNOWN"
}
