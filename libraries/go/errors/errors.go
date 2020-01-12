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

// GetMetadata returns the metadata map of the error
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
	case ErrPanic:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// Matches returns whether the string returned from error.Error() contains the given param string
func (e *Error) Matches(match string) bool {
	return strings.Contains(e.Error(), match)
}

// Generic error codes. Each of these has their own constructor for convenience.
const (
	ErrBadRequest         = "bad_request"
	ErrForbidden          = "forbidden"
	ErrInternalService    = "internal_service"
	ErrNotFound           = "not_found"
	ErrPreconditionFailed = "precondition_failed"
	ErrTimeout            = "timeout"
	ErrUnauthorized       = "unauthorized"
	ErrPanic              = "panic"
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

// WithCode wraps the error with a new code
func WithCode(err error, code string) *Error {
	return Wrap(err, code, "")
}

// WithMessage wraps the error with an extra message
func WithMessage(err error, format string, a ...interface{}) *Error {
	return Wrap(err, "", format, a...)
}

// WithMetadata will wrap the error with extra metadata
func WithMetadata(err error, metadata map[string]string) *Error {
	return Wrap(err, "", "", metadata)
}

// Wrap wraps the given error. Optionally, the last parameter can be a
// map[string]string containing metadata. If the error-to-wrap is
// already an *Error, the metadata will me merged and the existing
// code will remain the same. If the error is not an *Error, the
// code will default to ErrInternalService.
func Wrap(err interface{}, code, format string, a ...interface{}) *Error {
	// This allows the wrap functions to be called
	// without first checking whether the err is nil
	if err == nil {
		return nil
	}

	metadata, a := extractMetadata(format, a)

	// By default, the message of the returned error is the
	// error-to-wrap's message. If the given format is not
	// the empty string, the message becomes: new message: old message.
	var msg string
	switch v := err.(type) {
	case error:
		msg = v.Error()
	default:
		msg = fmt.Sprint(v)
	}

	if format != "" {
		msg = fmt.Sprintf(format, a...) + ": " + msg
	}

	// If the error to wrap is already an *Error
	switch v := err.(type) {
	case *Error:
		v.Message = msg
		v.Metadata = mergeMetadata(v.Metadata, metadata)

		if code != "" {
			v.Code = code
		}

		return v
	}

	if code == "" {
		code = ErrInternalService
	}

	return &Error{code, msg, metadata}
}

// newError returns a new Error with the given code. The message is formatted using Sprintf.
// If the last parameter is a map[string]string, it is assumed to be the error params.
func newError(code, format string, a []interface{}) *Error {
	metadata, a := extractMetadata(format, a)
	message := fmt.Sprintf(format, a...)
	return &Error{code, message, metadata}
}

func extractMetadata(format string, a []interface{}) (map[string]string, []interface{}) {
	if len(a) > 0 {
		// If we have too many parameters for the formatting directive,
		// the last parameter should be a metadata map.
		operandCount := countFmtOperands(format)
		if len(a) > operandCount {
			metadata, ok := a[len(a)-1].(map[string]string)
			if !ok {
				panic("Failed to assert metadata type")
			}
			return metadata, a[:operandCount]
		}
	}

	return nil, a
}

// mergeMetadata merges the metadata but preserves existing entries
func mergeMetadata(current, new map[string]string) map[string]string {
	if len(new) == 0 {
		return current
	}

	if current == nil {
		current = map[string]string{}
	}

	for k, v := range new {
		if _, ok := current[k]; !ok {
			current[k] = v
		}
	}

	return current
}
