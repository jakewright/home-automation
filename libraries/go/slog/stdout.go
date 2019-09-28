package slog

import "fmt"

// StdoutLogger writes all logs to stdout
type StdoutLogger struct{}

// NewStdoutLogger returns a StdoutLogger for the service with the given name
func NewStdoutLogger() Logger {
	return &StdoutLogger{}
}

// Log prints the event to stdout
func (l *StdoutLogger) Log(event *Event) {
	fmt.Println(event.String())
}
