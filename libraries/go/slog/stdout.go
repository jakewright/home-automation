package slog

import "fmt"

// Stdout logger writes all logs to stdout
type StdoutLogger struct{}

// NewStdoutLogger returns a StdoutLogger for the service with the given name
func NewStdoutLogger() Logger {
	return &StdoutLogger{}
}

func (l *StdoutLogger) Log(log *Log) {
	fmt.Println(log.String())
}
