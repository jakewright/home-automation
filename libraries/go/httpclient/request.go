package httpclient

import (
	"context"
	"io"
	"net/http"
)

// Request holds the information needed to make an HTTP request
type Request struct {
	Ctx             context.Context
	Method          string
	URL             string
	Headers         http.Header
	Body            interface{}
	RequestEncoder  Encoder
	ResponseDecoder Decoder
}

func (r *Request) validate() error {
	switch {
	case !validMethod(r.Method):
		return InvalidMethodError(r.Method)
	}
	return nil
}

func (r *Request) build(defaultEncoder Encoder) (*http.Request, error) {
	body, contentType, err := r.prepareBody(defaultEncoder)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(r.Method, r.URL, body)
	if err != nil {
		return nil, err
	}

	if r.Ctx != nil {
		req = req.WithContext(r.Ctx)
	}

	if r.Headers != nil {
		req.Header = r.Headers
	}

	// Set the Content-Type header (unless an override was provided in request)
	if contentType != "" && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", contentType)
	}

	return req, nil
}

func (r *Request) prepareBody(defaultEncoder Encoder) (io.Reader, string, error) {
	if r.Body == nil {
		return nil, "", nil
	}

	enc := r.RequestEncoder
	if enc == nil {
		enc = defaultEncoder
	}

	reader, err := enc.Encode(r.Body)
	if err != nil {
		return nil, "", err
	}

	return reader, enc.ContentType(), nil
}

func validMethod(method string) bool {
	switch method {
	case http.MethodGet:
	case http.MethodHead:
	case http.MethodPost:
	case http.MethodPut:
	case http.MethodPatch:
	case http.MethodDelete:
	case http.MethodConnect:
	case http.MethodOptions:
	case http.MethodTrace:
	default:
		return false
	}

	return true
}
