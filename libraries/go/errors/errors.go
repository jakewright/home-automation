package errors

import (
	"fmt"
	"net/http"
	"strings"
)

// Error is a custom error type that implements Go's error interface
type Error struct {
	Code     string            `json:"code"`
	Message  string            `json:"message"`
	Metadata map[string]string `json:"metadata"`
}

// Error returns a string message of the error
func (e *Error) Error() string {
	switch {
	case e == nil:
		return ""
	case e.Message == "":
		return e.Code
	case e.Code == "":
		return e.Message
	default:
		return fmt.Sprintf("%s: %s", e.Code, e.Message)
	}
}

func (e *Error) GetMetadata() map[string]string {
	return e.Metadata
}

// HTTPStatus returns an appropriate HTTP status code to use when returning the error in a response
func (e *Error) HTTPStatus() int {
	switch e.Code {
	case ErrBadRequest:
		return http.StatusBadRequest
	case ErrForbidden:
		return http.StatusForbidden
	case ErrInternalService:
		return http.StatusInternalServerError
	case ErrNotFound:
		return http.StatusNotFound
	case ErrPreconditionFailed:
		return http.StatusPreconditionFailed
	case ErrTimeout:
		return http.StatusRequestTimeout
	case ErrUnauthorized:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

// Matches returns whether the string returned from error.Error() contains the given param string
func (e *Error) Matches(match string) bool {
	return strings.Contains(e.Error(), match)
}

// Generic error codes. Each of these has their own constructor for convenience.
// You can use any string as a code, just use the `New` method.
const (
	ErrBadRequest         = "bad_request"
	ErrForbidden          = "forbidden"
	ErrInternalService    = "internal_service"
	ErrNotFound           = "not_found"
	ErrPreconditionFailed = "precondition_failed"
	ErrTimeout            = "timeout"
	ErrUnauthorized       = "unauthorized"
)

// InternalService creates a new error to represent an internal service error
func InternalService(format string, a ...interface{}) *Error {
	return newError(ErrInternalService, format, a)
}

// BadRequest creates a new error to represent an error caused by the client sending
// an invalid request. This is non-retryable unless the request is modified.
func BadRequest(format string, a ...interface{}) *Error {
	return newError(ErrBadRequest, format, a)
}

// Forbidden creates a new error representing a resource that cannot be accessed with
// the current authorisation credentials. The user may need authorising, or if authorised,
// may not be permitted to perform this action.
func Forbidden(format string, a ...interface{}) *Error {
	return newError(ErrForbidden, format, a)
}

// NotFound creates a new error representing a resource that cannot be found
func NotFound(format string, a ...interface{}) *Error {
	return newError(ErrNotFound, format, a)
}

// PreconditionFailed creates a new error indicating that one or more conditions
// given in the request evaluated to false when tested on the server
func PreconditionFailed(format string, a ...interface{}) *Error {
	return newError(ErrPreconditionFailed, format, a)
}

// Timeout creates a new error representing a timeout from client to server
func Timeout(format string, a ...interface{}) *Error {
	return newError(ErrTimeout, format, a)
}

// Unauthorized creates a new error indicating that authentication is required,
// but has either failed or not been provided.
func Unauthorized(format string, a ...interface{}) *Error {
	return newError(ErrUnauthorized, format, a)
}

func Wrap(err error, metadata map[string]string) *Error {
	return &Error{ErrInternalService, err.Error(), metadata}
}

// newError returns a new Error with the given code. The message is formatted using Sprintf.
// If the last parameter is a map[string]string, it is assumed to be the error params.
func newError(code, format string, params ...interface{}) *Error {
	// Take the last parameter
	last := params[len(params)-1]

	// Try to cast it to a map[string]string. If it fails, metadata will be an empty map.
	metadata, ok := last.(map[string]string)

	var message string

	// If the last parameter was a map[string]string
	if ok {
		// Format the string using all but the last parameter
		message = fmt.Sprintf(format, params[:len(params)-1]...)
	} else {
		// Format the string using all parameters
		message = fmt.Sprintf(format, params...)
	}

	return &Error{code, message, metadata}
}
