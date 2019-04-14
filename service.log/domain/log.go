package domain

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
)

const jsonIndent = "    "

type Line struct {
	Summary template.HTML
	Raw     template.HTML
}

type Log struct {
	Lines []*Line
}

func NewLineFromBytes(b []byte) *Line {
	return &Line{
		Summary: template.HTML(generateSummary(b)),
		Raw:     template.HTML(formatRaw(b)),
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

func generateSummary(b []byte) string {
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return string(b)
	}

	if _, ok := m["@timestamp"]; !ok {
		return string(b)
	}

	if _, ok := m["message"]; !ok {
		return string(b)
	}

	return fmt.Sprintf("%s %s", m["@timestamp"], m["message"])
}
