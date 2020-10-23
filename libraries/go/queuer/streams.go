package queuer

import (
	"context"
	"time"
)

// Message represents a message published to or received
// from a Redis stream.
type Message struct {
	// ID is the Redis ID of the message. When publishing,
	// this can be set to "*" to instruct Redis to generate
	// an ID. This is usually what you want. If left blank,
	// it will automatically be set to "*".
	ID string

	// Stream is the name of the stream this message should
	// be published to or was received from.
	Stream string

	// Values represents the message's payload
	Values map[string]interface{}

	retryCount int64
}

// Handler processes messages. HandleEvent should return a
// result that tells the client whether the handling was
// successful or not.
type Handler interface {
	HandleEvent(ctx context.Context, m *Message) Result
}

// HandlerFunc is an adapter that allows ordinary functions
// to be used as event handlers. If f is a function with the
// appropriate signature, HandlerFunc(f) is a Handler that
// calls f.
type HandlerFunc func(ctx context.Context, m *Message) Result

// HandleEvent calls f(e)
func (f HandlerFunc) HandleEvent(ctx context.Context, m *Message) Result {
	return f(ctx, m)
}

// Result defines the result of a handler
type Result struct {
	Retry   bool
	Err     error
	Backoff time.Duration
}

// Success should be returned when the message was
// successfully processed. The message will be acknowledged
// and therefore not re-processed by any consumers.
func Success() Result {
	return Result{}
}

// Fail should be returned when the message was not
// successfully processed and should be retried. The message
// will not be acknowledged and therefore retried. It will
// not necessarily be retried by the same consumer. The
// error will be enqueued on the error channel.
func Fail(err error) Result {
	return Result{
		Retry: true,
		Err:   err,
	}
}

// Discard should be returned when the message was not
// successfully processed but should not be retried. It will
// be acknowledged and therefore not re-processed by any
// consumers. The error will be enqueued on the error
// channel.
func Discard(err error) Result {
	return Result{
		Retry: false,
		Err:   err,
	}
}
