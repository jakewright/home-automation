package oops

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	pkgerrors "github.com/pkg/errors"
)

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

// Error is a custom error type that implements Go's error interface
type Error struct {
	code     string
	message  string
	metadata map[string]string
	cause    error
	stack    []uintptr
}

// GetCode unwinds the error stack and returns the
// first code encountered. If no codes exist, then
// ErrInternalService is returned.
func (e *Error) GetCode() string {
	if e == nil {
		return ""
	}

	if e.code != "" {
		return e.code
	}

	// If there's an inner error, return its
	// code provided it's not the empty string.
	if v, ok := e.cause.(*Error); ok {
		if code := v.GetCode(); code != "" {
			return code
		}
	}

	// Default to internal service
	return ErrInternalService
}

// GetMessage returns the error's message. If there
// is a cause, its error message is appended after
// a colon, and so on.
func (e *Error) GetMessage() string {
	var inner string
	if e.cause != nil {
		switch v := e.cause.(type) {
		case *Error:
			inner = v.GetMessage()
		default:
			inner = v.Error()
		}
	}

	return join(": ", e.message, inner)
}

// Error returns a string message of the error
func (e *Error) Error() string {
	if e == nil {
		return ""
	}

	// Call GetMessage() recursively instead of
	// Error() to avoid repeating the code
	return join(": ", e.GetCode(), e.GetMessage())
}

// GetMetadata returns the metadata map of the error
func (e *Error) GetMetadata() map[string]string {
	if e == nil {
		return nil
	}

	if v, ok := e.cause.(*Error); ok {
		return mergeMetadata(v.GetMetadata(), e.metadata)
	}

	return e.metadata
}

// HTTPStatus returns an appropriate HTTP status code to use when returning the error in a response
func (e *Error) HTTPStatus() int {
	switch e.GetCode() {
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

// StackTrace returns the stack of Frames from
// innermost (newest) to outermost (oldest)
func (e *Error) StackTrace() pkgerrors.StackTrace {
	f := make([]pkgerrors.Frame, len(e.stack))
	for i := 0; i < len(f); i++ {
		f[i] = pkgerrors.Frame(e.stack[i])
	}
	return f
}

// Matches returns whether the string returned from error.Error() contains the given param string
func (e *Error) Matches(match string) bool {
	// TODO make this useful
	return strings.Contains(e.Error(), match)
}

// InternalService creates a new error to represent an internal service error
func InternalService(format string, a ...interface{}) *Error {
	return newError(ErrInternalService, format, a, nil)
}

// BadRequest creates a new error to represent an error caused by the client sending
// an invalid request. This is non-retryable unless the request is modified.
func BadRequest(format string, a ...interface{}) *Error {
	return newError(ErrBadRequest, format, a, nil)
}

// Forbidden creates a new error representing a resource that cannot be accessed with
// the current authorisation credentials. The user may need authorising, or if authorised,
// may not be permitted to perform this action.
func Forbidden(format string, a ...interface{}) *Error {
	return newError(ErrForbidden, format, a, nil)
}

// NotFound creates a new error representing a resource that cannot be found
func NotFound(format string, a ...interface{}) *Error {
	return newError(ErrNotFound, format, a, nil)
}

// PreconditionFailed creates a new error indicating that one or more conditions
// given in the request evaluated to false when tested on the server
func PreconditionFailed(format string, a ...interface{}) *Error {
	return newError(ErrPreconditionFailed, format, a, nil)
}

// Timeout creates a new error representing a timeout from client to server
func Timeout(format string, a ...interface{}) *Error {
	return newError(ErrTimeout, format, a, nil)
}

// Unauthorized creates a new error indicating that authentication is required,
// but has either failed or not been provided.
func Unauthorized(format string, a ...interface{}) *Error {
	return newError(ErrUnauthorized, format, a, nil)
}

// Is returns whether the code matches that of the error
func Is(err error, code string) bool {
	if v, ok := err.(*Error); ok {
		return v.GetCode() == code
	}

	return false
}

// WithCode wraps the error with a new code
func WithCode(err interface{}, code string) *Error {
	return Wrap(err, code, "")
}

// WithMessage wraps the error with an extra message
func WithMessage(err interface{}, format string, a ...interface{}) *Error {
	return Wrap(err, "", format, a...)
}

// WithMetadata will wrap the error with extra metadata
func WithMetadata(err interface{}, metadata map[string]string) *Error {
	return Wrap(err, "", "", metadata)
}

// Wrap wraps the given error. Optionally, the last parameter can be a
// map[string]string containing metadata. If the error-to-wrap is
// already an *Error, the metadata will me merged and the existing
// code will remain the same. If the error is not an *Error, the
// code will default to ErrInternalService.
func Wrap(err interface{}, code, format string, a ...interface{}) *Error {
	// Accepting an interface allows us to wrap
	// things like the result of recover().
	var cause error
	switch v := err.(type) {
	case error:
		cause = v
	default:
		cause = errors.New(fmt.Sprint(v))
	}

	return newError(code, format, a, cause)
}

// newError returns a new Error with the given code. The message is formatted using Sprintf.
// If the last parameter is a map[string]string, it is assumed to be the error params.
func newError(code, format string, a []interface{}, cause error) *Error {
	metadata, a := extractMetadata(format, a)
	message := fmt.Sprintf(format, a...)
	return &Error{code, message, metadata, cause, stack()}
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

// join concatenates the elements of parts, placing sep
// in between each element. Empty strings are ignored.
func join(sep string, parts ...string) string {
	var str string
	for i, p := range parts {
		if p == "" {
			continue
		}

		if i > 0 {
			str += sep
		}

		str += p
	}
	return str
}

func stack() []uintptr {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(0, pcs[:])
	return pcs[0:n]
}
