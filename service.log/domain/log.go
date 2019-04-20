package domain

import (
	"bytes"
	"encoding/json"
	"home-automation/libraries/go/slog"
	"html/template"
	"time"
)

const jsonIndent = "    "

type Line struct {
	ID        int
	Timestamp time.Time     `json:"@timestamp"`
	Severity  slog.Severity `json:"severity"`
	Service   string        `json:"service"`
	Message   string        `json:"message"`
	Metadata  interface{}   `json:"metadata"`
	Raw       []byte        `json:"-"`
}

type FormattedLine struct {
	ID             int
	Timestamp      string
	Severity       string
	Service        string
	Message        template.HTML
	Metadata       template.HTML
	MetadataPretty template.HTML
	Raw            template.HTML
}

type Log struct {
	FormattedLines []*FormattedLine
}

func NewLineFromBytes(id int, b []byte) *Line {
	// Set some defaults in case they can't be parsed from the log
	l := Line{
		ID:       id,
		Severity: slog.InfoSeverity,
		Message:  string(b),
		Raw:      b,
	}

	if err := json.Unmarshal(b, &l); err != nil {
		// Ignore errors because there's no guarantee it's even JSON
		slog.Warn("Failed to unmarshal log line: %v", err)
	}

	if l.Timestamp.IsZero() {
		slog.Error("Log timestamp was zero: %v", string(l.Raw))
	}

	return &l
}

func (l *Line) Format() *FormattedLine {
	var metadataPretty []byte
	metadata, err := json.Marshal(l.Metadata)
	if err == nil {
		var buf bytes.Buffer
		err := json.Indent(&buf, metadata, "", jsonIndent)
		if err == nil {
			metadataPretty = buf.Bytes()
		}
	}

	raw := template.HTML(formatRaw(l.Raw))

	return &FormattedLine{
		ID:             l.ID,
		Timestamp:      l.Timestamp.Format(time.Stamp),
		Severity:       string(l.Severity),
		Service:        l.Service,
		Message:        template.HTML(l.Message),
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
