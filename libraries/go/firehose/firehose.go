package firehose

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"

	"github.com/jakewright/home-automation/libraries/go/oops"
	"github.com/jakewright/home-automation/libraries/go/queuer"
	"github.com/jakewright/home-automation/libraries/go/slog"
)

const firehoseStream = "firehose"

// Result is an alias for queuer.Result
type Result queuer.Result

// Success should be returned when the event was
// successfully processed.
func Success() Result {
	return Result{}
}

// Fail should be returned when the event was not
// successfully processed and should be retried
func Fail(err error) Result {
	return Result{
		Retry: true,
		Err:   err,
	}
}

// Discard should be returned when the event was not
// successfully processed but should not be retried
func Discard(err error) Result {
	return Result{
		Retry: false,
		Err:   err,
	}
}

// Decoder is a function that decodes an event into v
type Decoder func(v interface{}) error

// Handler processes messages received on a channel that has
// been subscribed to. HandleEvent should return a result
// that tells the client whether the handling was successful
// or not.
type Handler interface {
	HandleEvent(ctx context.Context, decode Decoder) Result
}

// HandlerFunc is an adapter that allows ordinary functions
// to be used as event handlers. If f is a function with the
// appropriate signature, HandlerFunc(f) is a Handler that
// calls f.
type HandlerFunc func(ctx context.Context, decode Decoder) Result

// HandleEvent calls f(ctx, decode)
func (f HandlerFunc) HandleEvent(ctx context.Context, decode Decoder) Result {
	return f(ctx, decode)
}

// Publisher is the interface that wraps the publish method
type Publisher interface {
	// Publish emits an event to the firehose
	Publish(ctx context.Context, channel string, message interface{}) error
}

// Subscriber is the interface that wraps the subscribe method
type Subscriber interface {
	Subscribe(channel string, handler Handler)
}

// ClientOptions contains fields to configure the Firehose client
type ClientOptions struct {
	Group    string
	Consumer string
	Redis    *redis.Client

	// HandlerTimeout is the duration after which the
	// context passed to handlers is cancelled. If this
	// value is zero, a default duration of 30s is applied.
	HandlerTimeout time.Duration
}

// Client implements the Publish and Subscribe interfaces
type Client struct {
	group    string
	consumer string

	c *queuer.Consumer
	p *queuer.Publisher

	handlers map[string]Handler
	mu       sync.Mutex
}

var _ Publisher = (*Client)(nil)
var _ Subscriber = (*Client)(nil)

// NewClient returns a new Firehose client
func NewClient(opts *ClientOptions) *Client {
	handlerTimeout := opts.HandlerTimeout
	if handlerTimeout == 0 {
		handlerTimeout = time.Second * 30
	}

	consumer := queuer.NewConsumer(&queuer.ConsumerOptions{
		Group:          opts.Group,
		Consumer:       opts.Consumer,
		Redis:          opts.Redis,
		ReadTimeout:    time.Second * 10,
		HandlerTimeout: handlerTimeout,
		PendingTimeout: handlerTimeout * 2,
		ClaimInterval:  handlerTimeout * 2,
		MaxRetry:       10,
		Concurrency:    3,
		BufferSize:     100,
		Backoff:        nil, // Use default
		NetworkRetry:   10,
	})

	publisher := queuer.NewPublisher(&queuer.PublisherOptions{
		StreamMaxLength:      1000,
		ApproximateMaxLength: true,
		Redis:                opts.Redis,
		Backoff:              nil, // Use default
		NetworkRetry:         10,
	})

	return &Client{
		group:    opts.Group,
		consumer: opts.Consumer,
		c:        consumer,
		p:        publisher,
		handlers: make(map[string]Handler),
	}
}

// Publish emits a message to the firehose
func (c *Client) Publish(ctx context.Context, channel string, message interface{}) error {
	data, err := json.Marshal(message)
	if err != nil {
		return oops.WithMessage(err, "failed to marshal payload")
	}

	if err := c.p.Publish(ctx, &queuer.Message{
		Stream: firehoseStream,
		Values: map[string]interface{}{
			"group":    c.group,
			"consumer": c.consumer,
			"channel":  channel,
			"data":     data,
		},
	}); err != nil {
		return oops.WithMessage(err, "failed to publish event")
	}

	return nil
}

// Subscribe registers the handler against the channel. If
// a handler already exists, the function will panic.
func (c *Client) Subscribe(channel string, handler Handler) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, ok := c.handlers[channel]; ok {
		slog.Panicf("already subscribed to Firehose channel %q", channel)
	}

	c.handlers[channel] = handler
}

// GetName returns a friendly name for the process
func (c *Client) GetName() string {
	return "Firehose"
}

// Start listens for messages on the firehose
func (c *Client) Start(ctx context.Context) error {
	var handler = func(ctx context.Context, m *queuer.Message) queuer.Result {
		channel, ok := m.Values["channel"].(string)
		if !ok {
			return queuer.Discard(oops.InternalService(
				"malformed Firehose event: no channel",
			))
		}

		c.mu.Lock()
		defer c.mu.Unlock()

		h, ok := c.handlers[channel]
		if !ok {
			return queuer.Success()
		}

		data, ok := m.Values["data"].(string)
		if !ok {
			return queuer.Discard(oops.InternalService(
				"malformed Firehose event: no data field",
			))
		}

		decode := func(v interface{}) error {
			if err := json.Unmarshal([]byte(data), &v); err != nil {
				return oops.WithMessage(err, "failed to decode Firehose event")
			}

			return nil
		}

		return queuer.Result(h.HandleEvent(ctx, decode))
	}

	var deadLetterHandler = func(ctx context.Context, m *queuer.Message) queuer.Result {
		slog.Errorf("dropping Firehose message %q", m.ID, map[string]string{
			"channel": m.Values["channel"].(string),
			"id":      m.ID,
		})
		return queuer.Success()
	}

	c.c.Subscribe(
		firehoseStream,
		queuer.HandlerFunc(handler),
		queuer.HandlerFunc(deadLetterHandler),
	)

	go func() {
		// The channel should be closed when Listen returns
		for err := range c.c.Errors {
			slog.Errorf("firehose error: %v", err)
		}
	}()

	if err := c.c.Listen(ctx); err != nil {
		return oops.WithMessage(err, "error returned from Listen()")
	}

	return nil
}
