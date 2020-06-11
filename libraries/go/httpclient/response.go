package httpclient

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
)

// Response represents the response from a request
type Response struct {
	*http.Response
}

// BodyBytes returns the body as a byte slice
func (r *Response) BodyBytes() ([]byte, error) {
	switch rc := r.Body.(type) {
	case *bufCloser:
		return rc.Bytes(), nil

	default:
		defer func() { _ = rc.Close() }()

		// Replace the response body with a bufCloser
		buf := &bufCloser{}
		r.Body = buf

		// Use a TeeReader to read the body while
		// simultaneously piping it into the buffer
		tr := io.TeeReader(rc, buf)
		return ioutil.ReadAll(tr)
	}
}

// BodyString returns the body as a string
func (r *Response) BodyString() (string, error) {
	b, err := r.BodyBytes()
	return string(b), err
}

type bufCloser struct {
	bytes.Buffer
}

// Close is a no-op
func (b *bufCloser) Close() error {
	return nil
}
