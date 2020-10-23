package queuer

import (
	"fmt"
	"net"
	"time"
)

// NetworkError is returned if a network error occurs when
// communicating with Redis. The Retrying field tells you
// whether the library is planning to retry the request. If
// this is true, Backoff will be set.
type NetworkError struct {
	Err      net.Error
	Retrying bool
	Backoff  time.Duration
}

// Unwrap returns the underlying error
func (e *NetworkError) Unwrap() error { return e.Err }

// Error returns a formatted error string
func (e *NetworkError) Error() string {
	if e.Retrying {
		return fmt.Sprintf(
			"failed to read Redis stream [retrying in %s]: %s",
			e.Backoff.String(), e.Err,
		)
	}

	return fmt.Sprintf(
		"failed to read Redis stream: %s",
		e.Err,
	)
}

// RedisError is returned if the Redis client returns an
// error that is not a network error.
type RedisError struct {
	Err error
}

// Unwrap returns the underlying error
func (e *RedisError) Unwrap() error { return e.Err }

// Error returns a formatted error string
func (e *RedisError) Error() string {
	return fmt.Sprintf("failed to read Redis stream: %s", e.Err)
}

// HandlerError wraps any errors returned from message handlers
type HandlerError struct {
	Err error
	Msg *Message
}

// Unwrap returns the underlying error
func (e *HandlerError) Unwrap() error { return e.Err }

// Error returns a formatted error string
func (e *HandlerError) Error() string {
	return fmt.Sprintf("error when handling message %q on stream %q: %s", e.Msg.ID, e.Msg.Stream, e.Err)
}

// HandlerPanic is returned if a message handler panics. The
// panic is recovered and converted to an error.
type HandlerPanic struct {
	Err error
	Msg *Message
}

// Unwrap returns the underlying error
func (e *HandlerPanic) Unwrap() error { return e.Err }

// Error returns a formatted error string
func (e *HandlerPanic) Error() string {
	return fmt.Sprintf("panic when handling message %q on stream %q: %s", e.Msg.ID, e.Msg.Stream, e.Err)
}
