package firehose

import (
	"context"
	"encoding/json"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis/v7"

	"github.com/jakewright/home-automation/libraries/go/oops"
	"github.com/jakewright/home-automation/libraries/go/slog"
)

// RedisClient wraps a redis.Client and exposes a Publish method
type RedisClient struct {
	MaxRetries int
	Backoff    time.Duration

	client *redis.Client
	pubsub *redis.PubSub

	handlers  map[string]func(*Event)
	phandlers map[string]func(*Event)

	shutdownInvoked *int32
	mux             sync.RWMutex
}

// NewRedisClient returns a RedisClient
func NewRedisClient(client *redis.Client) *RedisClient {
	return &RedisClient{
		MaxRetries: 3,
		Backoff:    time.Second * 3,

		client:          client,
		handlers:        make(map[string]func(*Event)),
		phandlers:       make(map[string]func(*Event)),
		shutdownInvoked: new(int32),
		mux:             sync.RWMutex{},
	}
}

// GetName returns a friendly name for the process
func (c *RedisClient) GetName() string {
	return "firehose"
}

// Publish emits the given message on the given redis channel
func (c *RedisClient) Publish(channel string, message interface{}) error {
	// The go-redis Publish() function will error if the message is not a string,
	// []byte, bool, or numeric type. Strings will be converted to []byte.

	data, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return c.client.Publish(channel, data).Err()
}

// Subscribe registers a handler for the given channel. Note that the subscription
// is not made with the Redis client until Start() is called.
func (c *RedisClient) Subscribe(channel string, handler Handler) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if _, ok := c.handlers[channel]; ok {
		slog.Panicf("Multiple handlers subscribed to the same channel")
	}

	wrappedHandler := func(e *Event) {
		params := map[string]string{
			"channel": e.Channel,
			"pattern": e.Pattern,
		}

		// Count the number of attempts
		e.attempts++
		var res Result

		for ; e.attempts <= c.MaxRetries; e.attempts++ {
			// Dispatch to the handler
			res = handler.HandleEvent(e)

			// If reached the maximum number of retries
			if e.attempts == c.MaxRetries {
				// This break means e.attempts is
				// not incremented beyond MaxRetries
				break
			}

			action := "discarding"
			if res.retry {
				action = "retrying"
			}

			if res.err != nil {
				slog.Errorf("Failed to handle event [attempt %d, %s...]: %v", e.attempts, action, res.err, params)
			}

			if !res.retry {
				return
			}

			// Back off before trying again
			time.Sleep(c.Backoff)
		}

		slog.Errorf("Failed to handle event [attempt %d, final]: %v", e.attempts, res.err, params)
	}

	c.handlers[channel] = wrappedHandler
}

// Start subscribes to the channels and listens for messages.
// The pubsub client is closed when the context is cancelled.
// Note that the underlying Redis client is *not* closed by
// this function. That should be handled by whatever passed
// in the Redis client to the New function in the first place.
func (c *RedisClient) Start(ctx context.Context) error {
	ch := make(chan error)

	go func() {
		ch <- c.listen()
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		if err := c.close(); err != nil {
			return err
		}
	}

	return <-ch
}

// listen subscribes to the channels and listens for messages
func (c *RedisClient) listen() error {
	c.pubsub = c.client.Subscribe()

	// Subscribe to channels
	var channels []string
	for channel := range c.handlers {
		channels = append(channels, channel)
	}
	if err := c.pubsub.Subscribe(channels...); err != nil {
		return err
	}

	// Subscribe to patterns
	var patterns []string
	for pattern := range c.phandlers {
		patterns = append(patterns, pattern)
	}
	if err := c.pubsub.PSubscribe(patterns...); err != nil {
		return err
	}

	// Wait for confirmation that all subscriptions have succeeded
	for subs := 0; subs < len(channels)+len(patterns); {
		msg, err := c.pubsub.ReceiveTimeout(time.Second * 10)
		if err != nil {
			return oops.WithMessage(err, "Timeout while waiting for Redis subscription confirmation")
		}

		switch v := msg.(type) {
		case *redis.Message, *redis.Pong:
			// Ignore
		case *redis.Subscription:
			subs++
			switch v.Kind {
			case "subscribe":
				slog.Infof("Subscribed to Redis channel %s", v.Channel)
			case "psubscribe":
				slog.Infof("Subscribed to Redis pattern %s", v.Channel)
			case "unsubscribe":
				return oops.InternalService("Unexpectedly unsubscribed from Redis channel %s", v.Channel)
			case "punsubscribe":
				return oops.InternalService("Unexpectedly unsubscribed from Redis pattern %s", v.Channel)
			}
		default:
			return oops.InternalService("Received unexpected message from Redis")
		}

		// Exit early if Stop() has been called
		if atomic.LoadInt32(c.shutdownInvoked) > 0 {
			return nil
		}
	}

	wg := sync.WaitGroup{}

	for msg := range c.pubsub.Channel() {
		params := map[string]string{
			"channel": msg.Channel,
			"pattern": msg.Pattern,
		}

		slog.Debugf("Received Redis message", params)

		event := Event{
			Channel: msg.Channel,
			Pattern: msg.Pattern,
			Payload: []byte(msg.Payload),
		}

		c.mux.RLock()

		// If there's a handler for this channel
		if handler, ok := c.handlers[msg.Channel]; ok {
			slog.Debugf("Dispatching Redis message to handler", params)

			wg.Add(1)
			go func(e Event) {
				handler(&e)
				wg.Done()
			}(event)
		}

		// If there's a handler for this pattern
		if handler, ok := c.phandlers[msg.Pattern]; ok {
			slog.Debugf("Dispatching Redis message to phandler", params)

			wg.Add(1)
			go func(e Event) {
				handler(&e)
				wg.Done()
			}(event)
		}

		c.mux.RUnlock()
	}

	// Wait for all of the handlers to finish
	wg.Wait()
	return nil
}

// close closes the pubsub channel so that the consumer
// stops receiving new messages. Start() will end once
// all in-flight handlers return.
func (c *RedisClient) close() error {
	// This is used to stop the listen() function
	// if it's stuck waiting for subscription confirmations
	atomic.StoreInt32(c.shutdownInvoked, 1)

	if c.pubsub == nil {
		return nil
	}

	return c.pubsub.Close()
}
