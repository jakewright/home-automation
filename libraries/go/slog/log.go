package slog

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Logger interface {
	Log(Severity, string, ...interface{})
}

type StdoutLogger struct {
	Service string
}

type metadataProvider interface {
	GetMetadata() map[string]string
}

type Severity string

const (
	DebugSeverity Severity = "debug"
	InfoSeverity  Severity = "info"
	WarnSeverity  Severity = "warning"
	ErrorSeverity Severity = "err"
)

type Log struct {
	Timestamp time.Time         `json:"timestamp"`
	Service   string            `json:"service"`
	Severity  Severity          `json:"severity"`
	Message   string            `json:"message"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

var DefaultLogger Logger

func mustGetDefaultLogger() Logger {
	if DefaultLogger == nil {
		panic("Default logger not set")
	}

	return DefaultLogger
}

func New(service string) Logger {
	return &StdoutLogger{
		Service: service,
	}
}

func Debug(format string, params ...interface{}) {
	mustGetDefaultLogger().Log(DebugSeverity, format, params...)
}

func Info(format string, params ...interface{}) {
	mustGetDefaultLogger().Log(InfoSeverity, format, params...)
}

func Warn(format string, params ...interface{}) {
	mustGetDefaultLogger().Log(WarnSeverity, format, params...)
}

func Error(format string, params ...interface{}) {
	mustGetDefaultLogger().Log(ErrorSeverity, format, params...)
}

func Fatal(format string, params ...interface{}) {
	Error(format, params...)
	os.Exit(1)
}

func (l *StdoutLogger) Log(severity Severity, format string, params ...interface{}) {
	// Take the last parameter
	last := params[len(params)-1]

	// Try to cast it to a map[string]string. If it fails, metadata will be an empty map.
	metadata, ok := last.(map[string]string)

	var message string

	// If the last parameter was a map[string]string
	if ok {
		// Format the string using all but the last parameter
		message = fmt.Sprintf(format, params[:len(params)-1]...)
	} else {
		// Format the string using all parameters
		message = fmt.Sprintf(format, params...)
	}

	// If any of the parameters have their own metadata (e.g. an Error),
	// merge it with the existing metadata.
	for _, param := range params {
		if param, ok := param.(metadataProvider); ok {
			metadata = mergeMetadata(metadata, param.GetMetadata())
		}
	}

	l.log(&Log{
		Timestamp: time.Now(),
		Service:   l.Service,
		Severity:  Severity(severity),
		Message:   message,
		Metadata:  metadata,
	})
}

func (l *StdoutLogger) log(log *Log) {
	b, err := json.Marshal(log)
	if err != nil {
		fmt.Println("Failed to marshal log line")
		return
	}

	fmt.Println(string(b))
}

// mergeMetadata merges the metadata but preserves existing entries
func mergeMetadata(current, new map[string]string) map[string]string {
	if len(new) == 0 {
		return current
	}

	if current == nil {
		current = map[string]string{}
	}

	for k, v := range new {
		if _, ok := current[k]; !ok {
			current[k] = v
		}
	}

	return current
}
