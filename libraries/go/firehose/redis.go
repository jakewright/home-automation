package firehose

import (
	"context"
	"encoding/json"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis/v7"

	"github.com/jakewright/home-automation/libraries/go/errors"
	"github.com/jakewright/home-automation/libraries/go/slog"
)

// RedisClient wraps a redis.Client and exposes a Publish method
type RedisClient struct {
	client *redis.Client
	pubsub *redis.PubSub
	cfg    *Config

	handlers  map[string]RawHandlerFunc
	phandlers map[string]RawHandlerFunc

	shutdownInvoked *int32
	mux             sync.RWMutex
}

// NewRedisClient returns a RedisClient
func NewRedisClient(client *redis.Client) *RedisClient {
	return &RedisClient{
		client:          client,
		handlers:        make(map[string]RawHandlerFunc),
		phandlers:       make(map[string]RawHandlerFunc),
		shutdownInvoked: new(int32),
		mux:             sync.RWMutex{},
	}
}

// WithConfig sets the client's config
func (c *RedisClient) WithConfig(cfg *Config) *RedisClient {
	c.cfg = cfg
	return c
}

func (c *RedisClient) config() *Config {
	if c.cfg == nil {
		return defaultConfig
	}

	return c.cfg
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
func (c *RedisClient) Subscribe(channel string, handler RawHandlerFunc) {
	c.mux.Lock()
	defer c.mux.Unlock()

	if _, ok := c.handlers[channel]; ok {
		slog.Panic("Multiple handlers subscribed to the same channel")
	}

	c.handlers[channel] = handler
}

// Start subscribes to the channels and listens for messages
func (c *RedisClient) Start() error {
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
			return errors.Wrap(err, "Timeout while waiting for Redis subscription confirmation")
		}

		switch v := msg.(type) {
		case *redis.Message, *redis.Pong:
			// Ignore
		case *redis.Subscription:
			subs++
			switch v.Kind {
			case "subscribe":
				slog.Info("Subscribed to Redis channel %s", v.Channel)
			case "psubscribe":
				slog.Info("Subscribed to Redis pattern %s", v.Channel)
			case "unsubscribe":
				return errors.InternalService("Unexpectedly unsubscribed from Redis channel %s", v.Channel)
			case "punsubscribe":
				return errors.InternalService("Unexpectedly unsubscribed from Redis pattern %s", v.Channel)
			}
		default:
			return errors.InternalService("Received unexpected message from Redis")
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

		slog.Debug("Received Redis message", params)

		event := Event{
			Channel: msg.Channel,
			Pattern: msg.Pattern,
			Payload: []byte(msg.Payload),
		}

		c.mux.RLock()

		// If there's a handler for this channel
		if handler, ok := c.handlers[msg.Channel]; ok {
			slog.Debug("Dispatching Redis message to handler", params)

			wg.Add(1)
			go func(e Event) {
				handler(&e)
				wg.Done()
			}(event)
		}

		// If there's a handler for this pattern
		if handler, ok := c.phandlers[msg.Pattern]; ok {
			slog.Debug("Dispatching Redis message to phandler", params)

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

// Stop closes the pubsub channel so that the consumer
// stops receiving new messages. Start() will end once
// all in-flight handlers return.
func (c *RedisClient) Stop(_ context.Context) error {
	// This is used to stop the Start() function
	// if it's stuck waiting for subscription confirmations
	atomic.StoreInt32(c.shutdownInvoked, 1)

	if c.pubsub == nil {
		return nil
	}

	return c.pubsub.Close()
}
