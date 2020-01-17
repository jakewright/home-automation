package slog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
)

type metadataProvider interface {
	GetMetadata() map[string]string
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// Event represents a single log event
type Event struct {
	Timestamp time.Time
	Severity  Severity
	Message   string
	Metadata  map[string]string
}

func newEvent(severity Severity, v interface{}) *Event {
	return newEventFromFormat(severity, "%v", v)
}

func newEventFromFormat(severity Severity, format string, a ...interface{}) *Event {
	metadata := make(map[string]string)

	if len(a) > 0 {
		// If we have too many parameters for the formatting directive,
		// the last parameter should be a metadata map.
		operandCount := countFmtOperands(format)
		if len(a) > operandCount {
			var ok bool
			metadata, ok = a[len(a)-1].(map[string]string)
			if !ok {
				panic("failed to assert metadata type")
			}
			a = a[:operandCount]
		}
	}

	message := fmt.Sprintf(format, a...)

	// If any of the parameters have their own metadata (e.g. an Error),
	// merge it with the existing metadata.
	for _, param := range a {
		if param, ok := param.(metadataProvider); ok {
			metadata = mergeMetadata(metadata, param.GetMetadata())
		}
		if param, ok := param.(stackTracer); ok {
			metadata["stack"] = fmt.Sprintf("%+v", param.StackTrace())
		}
	}

	return &Event{
		Timestamp: time.Now(),
		Severity:  severity,
		Message:   message,
		Metadata:  metadata,
	}
}

func (e *Event) String() string {
	metadata, err := json.Marshal(e.Metadata)
	if err != nil {
		fmt.Println("Failed to marshal metadata")
	}

	// If the JSON came out as "null" or "{}", don't bother writing anything.
	if bytes.Equal(metadata, []byte("null")) || bytes.Equal(metadata, []byte("{}")) {
		metadata = nil
	}

	timestamp := e.Timestamp.Format(time.RFC3339)
	return strings.Join([]string{timestamp, e.Severity.String(), e.Message, string(metadata)}, " ")
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
