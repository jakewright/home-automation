package firehose

// Handler processes messages received on a channel that has
// been subscribed to. HandleEvent should return a result
// that tells the client whether the handling was successful
// or not. The side-effects of returning non-successful
// results depend on the Subscriber implementation.
type Handler interface {
	HandleEvent(Event) Result
}

// HandlerFunc is an adapter that allows ordinary functions
// to be used as event handlers. If f is a function with the
// appropriate signature, HandlerFunc(f) is a Handler that
// calls f.
type HandlerFunc func(Event) Result

// HandleEvent calls f(e)
func (f HandlerFunc) HandleEvent(e Event) Result {
	return f(e)
}

// Publisher is the interface that wraps the publish method
type Publisher interface {
	// Publish emits an event to the firehose
	Publish(channel string, message interface{}) error
}

// Subscriber is the interface that wraps the subscribe method
type Subscriber interface {
	Subscribe(channel string, handler Handler)
}

// Event represents a message received from the Firehose
type Event interface {
	Channel() string
	Decode(v interface{}) error
}

// Result defines the result of a handler
type Result struct {
	retry bool
	err   error
}

// Success should be returned when the
// event was successfully processed
func Success() Result {
	return Result{}
}

// Fail should be returned when the event was not
// successfully processed and should be retried
func Fail(err error) Result {
	return Result{
		retry: true,
		err:   err,
	}
}

// Discard should be returned when the event was not
// successfully processed but should _not_ be retried
func Discard(err error) Result {
	return Result{
		retry: false,
		err:   err,
	}
}
