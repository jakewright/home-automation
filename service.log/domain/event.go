package domain

import (
	"bytes"
	"encoding/json"
	"html/template"
	"time"

	"github.com/jakewright/home-automation/libraries/go/slog"
)

const jsonIndent = "    "

// Event represents a single log event
type Event struct {
	// UUID is a unique identifier for this event
	UUID string `json:"uuid"`

	// Timestamp is the time on the event
	Timestamp time.Time `json:"@timestamp"`

	// Severity is the severity of the event
	Severity slog.Severity `json:"severity"`

	// Service is the name of the service from which the event came
	Service string `json:"service"`

	// Message is a best-effort attempt at parsing the message
	// from the raw line. This is handled by logstash.
	Message string `json:"message"`

	// Metadata is extra parameters put into slog lines.
	// This will usually be a map[string]string.
	Metadata interface{} `json:"metadata"`

	// Raw is the original log line
	Raw []byte `json:"-"`
}

// FormattedEvent is a version of Event that
// is ready to pass to an HTML template
type FormattedEvent struct {
	// UUID is a unique identifier for this event
	UUID string

	// Timestamp is the time on the event
	Timestamp string

	// Severity is the severity of the event
	Severity string

	// Service is the name of the service from which the event came
	Service string

	// Message is converted to template.HTML as-is
	Message template.HTML

	// Metadata is converted to template.HTML in its raw form
	Metadata template.HTML

	// MetadataPretty is parsed as JSON and indented if successful,
	// otherwise it will be equal to Metadata.
	MetadataPretty template.HTML

	// Raw is the original log line
	Raw template.HTML
}

// NewEventFromBytes returns a structured event from a log line.
// It is best-effort and therefore does not return an error. This
// is why this approach is used over a custom JSON unmarshal function.
func NewEventFromBytes(b []byte) *Event {
	// Set some defaults in case they can't be parsed from the log
	e := Event{
		Message: string(b),
		Raw:     b,
	}

	if err := json.Unmarshal(b, &e); err != nil {
		// Ignore errors because there's no guarantee it's even JSON
		slog.Warnf("Failed to unmarshal event: %v", err)
	}

	if e.Timestamp.IsZero() {
		slog.Warnf("Event timestamp was zero: %v", string(e.Raw))
	}

	return &e
}

// Format returns a formatted event that can be passed to an HTML template
func (e *Event) Format() *FormattedEvent {
	var metadataPretty []byte
	metadata, err := json.Marshal(e.Metadata)
	if err == nil {
		var buf bytes.Buffer
		err := json.Indent(&buf, metadata, "", jsonIndent)
		if err == nil {
			metadataPretty = buf.Bytes()
		}
	}

	raw := template.HTML(formatRaw(e.Raw))

	return &FormattedEvent{
		UUID:           e.UUID,
		Timestamp:      e.Timestamp.Format(time.Stamp),
		Severity:       e.Severity.String(),
		Service:        e.Service,
		Message:        template.HTML(e.Message),
		Metadata:       template.HTML(metadata),
		MetadataPretty: template.HTML(metadataPretty),
		Raw:            raw,
	}
}

func formatRaw(b []byte) string {
	var buf bytes.Buffer
	err := json.Indent(&buf, b, "", jsonIndent)
	if err != nil {
		return string(b)
	}

	return buf.String()
}
