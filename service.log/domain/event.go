package domain

import (
	"bytes"
	"encoding/json"
	"home-automation/libraries/go/slog"
	"html/template"
	"time"
)

const jsonIndent = "    "

type Event struct {
	ID        int
	Timestamp time.Time     `json:"@timestamp"`
	Severity  slog.Severity `json:"severity"`
	Service   string        `json:"service"`
	Message   string        `json:"message"`
	Metadata  interface{}   `json:"metadata"`
	Raw       []byte        `json:"-"`
}

type FormattedEvent struct {
	ID             int
	Timestamp      string
	Severity       string
	Service        string
	Message        template.HTML
	Metadata       template.HTML
	MetadataPretty template.HTML
	Raw            template.HTML
}

func NewEventFromBytes(id int, b []byte) *Event {
	// Set some defaults in case they can't be parsed from the log
	e := Event{
		ID: id,
		//Severity: slog.InfoSeverity,
		Message: string(b),
		Raw:     b,
	}

	if err := json.Unmarshal(b, &e); err != nil {
		// Ignore errors because there's no guarantee it's even JSON
		slog.Warn("Failed to unmarshal event: %v", err)
	}

	if e.Timestamp.IsZero() {
		slog.Error("Event timestamp was zero: %v", string(e.Raw))
	}

	return &e
}

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
		ID:             e.ID,
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
