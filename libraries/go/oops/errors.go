package oops

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"

	pkgerrors "github.com/pkg/errors"
)

// Code is an identifier of a particular class of error
type Code string

// Generic error codes. Each of these has their own constructor for convenience.
const (
	ErrBadRequest         Code = "bad_request"
	ErrForbidden          Code = "forbidden"
	ErrInternalService    Code = "internal_service"
	ErrNotFound           Code = "not_found"
	ErrPreconditionFailed Code = "precondition_failed"
	ErrTimeout            Code = "timeout"
	ErrUnauthorized       Code = "unauthorized"
)

var httpStatusByCode = map[Code]int{
	ErrBadRequest:         http.StatusBadRequest,
	ErrForbidden:          http.StatusForbidden,
	ErrInternalService:    http.StatusInternalServerError,
	ErrNotFound:           http.StatusNotFound,
	ErrPreconditionFailed: http.StatusPreconditionFailed,
	ErrTimeout:            http.StatusRequestTimeout,
	ErrUnauthorized:       http.StatusUnauthorized,
}

var codeByHTTPStatus map[int]Code

func init() {
	codeByHTTPStatus = make(map[int]Code, len(httpStatusByCode))
	for code, status := range httpStatusByCode {
		codeByHTTPStatus[status] = code
	}
}

// Error is a custom error type that implements Go's error interface
type Error struct {
	code     Code
	message  string
	metadata map[string]string
	cause    error
	stack    []uintptr
}

// GetCode unwinds the error stack and returns the
// first code encountered. If no codes exist, then
// ErrInternalService is returned.
func (e *Error) GetCode() Code {
	if e == nil {
		return ErrInternalService
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
	return join(": ", string(e.GetCode()), e.GetMessage())
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
	if status := httpStatusByCode[e.GetCode()]; status != 0 {
		return status
	}

	return http.StatusInternalServerError
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

// FromHTTPStatus returns an error where the code is derived from the HTTP status code
func FromHTTPStatus(status int, format string, a ...interface{}) *Error {
	code := codeByHTTPStatus[status]
	if code == "" {
		code = ErrInternalService
	}

	return newError(code, format, a, nil)
}

// Is returns whether the code matches that of the error
func Is(err error, code Code) bool {
	if v, ok := err.(*Error); ok {
		return v.GetCode() == code
	}

	return false
}

// WithCode wraps the error with a new code
func WithCode(err interface{}, code Code) *Error {
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
func Wrap(err interface{}, code Code, format string, a ...interface{}) *Error {
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
func newError(code Code, format string, a []interface{}, cause error) *Error {
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
	var written bool

	for _, p := range parts {
		if p == "" {
			continue
		}

		// Only print a separator if we've already written
		// something to avoid the case where the first n
		// parts being empty strings results in the returned
		// string beginning with the separator.
		if written {
			str += sep
		}

		str += p
		written = true
	}

	return str
}

func stack() []uintptr {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(0, pcs[:])
	return pcs[0:n]
}
