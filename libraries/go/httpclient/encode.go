package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Encoder is the interface for types that can encode a request body.
type Encoder interface {
	ContentType() string
	Encode(interface{}) (io.Reader, error)
}

// EncoderJSON encodes bodies as JSON
type EncoderJSON struct{}

// ContentType returns the ContentType header to set in an outbound request
func (e EncoderJSON) ContentType() string { return "application/json; charset=utf-8" }

// Encode marshals an arbitrary data structure into JSON
func (e EncoderJSON) Encode(body interface{}) (io.Reader, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to encode body as JSON: %w", err)
	}

	return bytes.NewReader(b), nil
}

// EncoderFormURL encodes bodies as x-www-form-urlencoded
type EncoderFormURL struct {
	// CustomContentType overrides the default ContentType
	// of application/x-www-form-urlencoded
	CustomContentType string

	// TagAlias is the tag to read on struct fields.
	// If empty, "form" will be used.
	TagAlias string
}

// ContentType returns the ContentType header to set in an outbound request
func (e EncoderFormURL) ContentType() string {
	if e.CustomContentType != "" {
		return e.CustomContentType
	}

	return "application/x-www-form-urlencoded"
}

// Encode marshals maps and simple structs to a URL-encoded body
func (e EncoderFormURL) Encode(body interface{}) (io.Reader, error) {
	tagAlias := "form"
	if e.TagAlias != "" {
		tagAlias = e.TagAlias
	}

	values, err := toQueryString(body, tagAlias)
	if err != nil {
		return nil, fmt.Errorf("failed to Encode body as URL query: %w", err)
	}

	return strings.NewReader(values), nil
}
