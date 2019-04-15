package slog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type metadataProvider interface {
	GetMetadata() map[string]string
}

// Log represents a single log event
type Log struct {
	Timestamp time.Time
	Severity  Severity
	Message   string
	Metadata  map[string]string
}

func newFromFormat(severity Severity, format string, params ...interface{}) *Log {
	// Take the last parameter
	var last interface{}
	if len(params) > 0 {
		last = params[len(params)-1]
	} else {
		last = nil
	}

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

	return &Log{
		Timestamp: time.Now(),
		Severity:  Severity(severity),
		Message:   message,
		Metadata:  metadata,
	}
}

func (l *Log) String() string {
	metadata, err := json.Marshal(l.Metadata)
	if err != nil {
		fmt.Println("Failed to marshal metadata")
	}

	// If the JSON came out as "null", don't bother writing anything.
	if bytes.Equal(metadata, []byte("null")) {
		metadata = nil
	}

	timestamp := l.Timestamp.Format(time.RFC3339)
	return strings.Join([]string{timestamp, string(l.Severity), l.Message, string(metadata)}, " ")
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
