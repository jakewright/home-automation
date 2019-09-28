package slog

import (
	"encoding/json"
	"strings"
)

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

	// UnknownSeverity is the value used when the severity cannot be derived
	UnknownSeverity Severity = 10
)

// String returns the name of the severity level
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

// UnmarshalJSON unmarshals a JSON string into a Severity
func (s *Severity) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	switch strings.ToLower(str) {
	case "dbg", "debug":
		*s = DebugSeverity
	case "inf", "info", "information":
		*s = InfoSeverity
	case "warn", "warning":
		*s = WarnSeverity
	case "err", "error":
		*s = ErrorSeverity
	default:
		*s = UnknownSeverity
	}

	return nil
}
