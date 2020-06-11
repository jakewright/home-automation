package httpclient

import (
	"encoding/json"
	"fmt"
)

// Decoder is the interface for types that can decode a response body.
type Decoder interface {
	Name() string
	Decode([]byte, interface{}) error
}

func inferDecoder(contentType string) (Decoder, error) {
	switch contentType {
	case "application/json",
		"application/json; charset=utf-8":
		return &DecoderJSON{}, nil
	}

	return nil, ContentTypeError(contentType)
}

// DecoderJSON decodes JSON bodies
type DecoderJSON struct{}

// Name returns the name of the format
func (d *DecoderJSON) Name() string {
	return "JSON"
}

// Decode unmarshals a JSON response body
func (d *DecoderJSON) Decode(data []byte, v interface{}) error {
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("failed to decode body as JSON: %w", err)
	}

	return nil
}
