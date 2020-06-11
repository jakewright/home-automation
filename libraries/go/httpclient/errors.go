package httpclient

import (
	"fmt"
	"net/http"
)

// InvalidMethodError is returned if an unsupported HTTP method is specified
type InvalidMethodError string

// Error implements the error interface
func (method InvalidMethodError) Error() string {
	return fmt.Sprintf("invalid method %q", method)
}

// BadStatusError is returned if the client's
// status validator function returns false.
type BadStatusError int

// Error implements the error interface
func (statusCode BadStatusError) Error() string {
	return fmt.Sprintf("request failed with status %d %s", statusCode, http.StatusText(int(statusCode)))
}

// ContentTypeError is returned if the response has a Content-Type
// header that cannot be handled by one of the built-in decoders and
// neither the client nor request have a default decoder set.
type ContentTypeError string

// Error implements the error interface
func (contentType ContentTypeError) Error() string {
	return fmt.Sprintf("unsupported Content-Type in response %q", contentType)
}

// DecodeError is a type alias used for errors during decoding
// so that they can be identified later by errors.As().
type DecodeError struct {
	Format string
	Err    error
}

// Error implements the error interface
func (e *DecodeError) Error() string {
	return fmt.Sprintf("failed to decode body as %s: %s", e.Format, e.Err)
}

// Unwrap returns the underlying error
func (e *DecodeError) Unwrap() error { return e.Err }
