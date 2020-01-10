package firehose

import (
	"time"

	"github.com/jakewright/home-automation/libraries/go/slog"
)

// Config defines how an event should be handled
type Config struct {
	MaxRetries int
	Backoff    time.Duration
}

var defaultConfig = &Config{
	MaxRetries: 3,
	Backoff:    time.Second * 3,
}

// Client is the interface that wraps the Publish and Subscribe methods
type Client interface {
	Publish(string, interface{}) error
	Subscribe(string, RawHandlerFunc)
	config() *Config
}

// DefaultPublisher is a global instance of Publisher
var DefaultClient Client

func mustGetDefaultClient() Client {
	if DefaultClient == nil {
		panic("Firehose used before default client set. Have you passed the Firehose option to bootstrap.Init()?")
	}

	return DefaultClient
}

// Publish sends the given message on the given channel using the default publisher
func Publish(channel string, message interface{}) error {
	return mustGetDefaultClient().Publish(channel, message)
}

// Subscribe offers syntactic sugar over the DefaultSubscriber's
// Subscribe function. Namely, it takes a HandlerFunc
func Subscribe(channel string, handler HandlerFunc) {
	c := mustGetDefaultClient()

	wrappedHandler := func(e *Event) {
		params := map[string]string{
			"channel": e.Channel,
			"pattern": e.Pattern,
		}

		// Count the number of attempts
		e.attempts++
		var res Result

		for ; e.attempts <= c.config().MaxRetries; e.attempts++ {
			// Dispatch to the handler
			res = handler(e)

			// If reached the maximum number of retries
			if e.attempts == c.config().MaxRetries {
				// This break means e.attempts is
				// not incremented beyond MaxRetries
				break
			}

			action := "discarding"
			if res.retry {
				action = "retrying"
			}

			if res.err != nil {
				slog.Error("Failed to handle event [attempt %d, %s...]: %v", e.attempts, action, res.err, params)
			}

			if !res.retry {
				return
			}

			// Back off before trying again
			time.Sleep(c.config().Backoff)
		}

		slog.Error("Failed to handle event [attempt %d, final]: %v", e.attempts, res.err, params)
	}

	c.Subscribe(channel, wrappedHandler)
}

// HandlerFunc is used by the syntactic sugar Subscribe() func
type HandlerFunc func(*Event) Result

// RawHandlerFunc is what a Subscriber dispatches messages to
type RawHandlerFunc func(*Event)

// Event represents a message received from the Firehose
type Event struct {
	Channel  string
	Pattern  string
	Payload  []byte
	attempts int
}

// Result defines the result of a handler
type Result struct {
	retry bool
	err   error
}

// Success should be returned when the event was successfully processed
func Success() Result {
	return Result{}
}

// Fail should be returned when the event was not successfully processed and should be retried
func Fail(err error) Result {
	return Result{
		retry: true,
		err:   err,
	}
}

// Discard should be returned when the event was not successfully processed but should not be retried
func Discard(err error) Result {
	return Result{
		retry: false,
		err:   err,
	}
}
