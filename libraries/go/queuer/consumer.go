package queuer

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jpillora/backoff"
	"golang.org/x/sync/errgroup"

	"github.com/go-redis/redis/v8"
)

// Redis errors
const (
	errConsumerGroupAlreadyExists = "BUSYGROUP Consumer Group name already exists"
)

// ConsumerOptions contains options to configure the Consumer
type ConsumerOptions struct {
	// Group is the name of the consumer group used when
	// listening for new messages from Redis. This must
	// be set before calling Listen().
	Group string

	// Consumer is the name of this particular consumer
	// within the consumer group. This must be set before
	// calling Listen().
	Consumer string

	// Redis is an instance of *redis.Client for use by
	// the client. This must be set before using the client.
	Redis *redis.Client

	// ReadTimeout is the duration for which the XREADGROUP
	// call blocks for. A duration of zero means the client
	// will block indefinitely. It is recommended to set
	// this to a non-zero duration so that the client is
	// able to gracefully shutdown.
	ReadTimeout time.Duration

	// HandlerTimeout is the duration after which the
	// context passed to handlers is cancelled. Note that
	// handlers are not forcefully stopped after this time.
	// It is up to them to handle context cancellation.
	// A duration of zero means handlers never timeout.
	HandlerTimeout time.Duration

	// PendingTimeout is the duration for which a message
	// can be pending before the consumer tries to claim it.
	// Note that a consumer can claim its own messages if
	// they are not acknowledged after processing. This is
	// how the retry mechanism works.
	//
	// This value should not be shorter than HandlerTimeout
	// otherwise you risk claiming messages that are still
	// being processed.
	PendingTimeout time.Duration

	// ClaimInterval is the time between attempts to claim
	// any messages that have been pending for longer than
	// the PendingTimeout. If this value is zero, then the
	// consumer will not try to claim pending messages. If
	// all consumers in the group have this set to zero then
	// unacknowledged messages will never be retried.
	ClaimInterval time.Duration

	// MaxRetry is the number of times a message will be
	// retried before it is passed to the dead-letter
	// consumer(s) for the stream. If < 0, then the
	// message will never be dead-lettered.
	MaxRetry int

	// Concurrency is the number of goroutines that are
	// spawned to concurrently handle incoming messages.
	// A value of zero is equal to a value of one.
	Concurrency int

	// BufferSize is the size of the channel that holds
	// incoming messages and therefore determines how many
	// messages the consumer can read from Redis in a
	// single call. A value of zero will create an
	// unbuffered channel.
	BufferSize int

	// Backoff is used to retry requests to Redis in the
	// case of network failures. If this value is nil, a
	// Backoff with sensible defaults will be used.
	Backoff *backoff.Backoff

	// NetworkRetry is the number of times to retry failed
	// network requests to Redis before returning a fatal
	// error. A value of zero means requests will not be
	// retried. A value of < 0 means requests will be
	// retried indefinitely until the context is cancelled.
	NetworkRetry int
}

// Consumer uses Redis Streams to publish and subscribe
type Consumer struct {
	// Errors is a channel over which non-fatal errors are
	// sent. This channel must have a listener otherwise
	// deadlock will arise.
	Errors chan error

	opts *ConsumerOptions

	handlers           map[string]Handler
	deadLetterHandlers map[string]Handler
	mu                 sync.Mutex // guards the maps

	messages chan *Message
	stopWork chan struct{}
	backoff  *backoff.Backoff
}

// NewConsumer returns an initialised Consumer
func NewConsumer(opts *ConsumerOptions) *Consumer {
	b := opts.Backoff
	if b == nil {
		b = defaultBackoff
	}

	return &Consumer{
		Errors:             make(chan error),
		opts:               opts,
		handlers:           make(map[string]Handler),
		deadLetterHandlers: make(map[string]Handler),
		mu:                 sync.Mutex{},
		messages:           make(chan *Message, opts.BufferSize),
		stopWork:           make(chan struct{}, opts.Concurrency),
		backoff:            b,
	}
}

// Subscribe registers a handler for the given channel
func (c *Consumer) Subscribe(channel string, handler Handler, deadLetter ...Handler) {
	c.mu.Lock()
	defer c.mu.Unlock()

	switch len(deadLetter) {
	case 0: // ignore
	case 1:
		c.deadLetterHandlers[channel] = deadLetter[0]
	default:
		panic(fmt.Errorf("too many dead letter handlers"))
	}

	c.handlers[channel] = handler
}

// Listen subscribes to the channels and listens for messages
func (c *Consumer) Listen(ctx context.Context) error {
	err := c.listen(ctx)
	switch {
	case errors.Is(err, context.DeadlineExceeded),
		errors.Is(err, context.Canceled):
		return nil
	}
	return err
}

func (c *Consumer) listen(ctx context.Context) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.handlers) == 0 {
		return nil
	}

	// Create a list of streams to listen to
	streams := make([]string, 0, len(c.handlers)*2)

	for stream := range c.handlers {
		streams = append(streams, stream)

		// Create the consumer group
		if err := c.xGroupCreateMkStream(ctx, stream); err != nil {
			return err
		}
	}

	for range streams {
		streams = append(streams, ">")
	}

	wg := &sync.WaitGroup{}

	// Start the workers
	for i := 0; i < c.opts.Concurrency; i++ {
		go func() {
			defer wg.Done()
			c.work(ctx)
		}()
	}

	// Start polling and reclaiming
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return c.poll(ctx, streams)
	})

	g.Go(func() error {
		return c.claim(ctx)
	})

	// Wait for polling and reclaiming to stop before
	// stopping the workers. If workers are stopped too
	// early, polling/reclaiming could deadlock trying
	// to write to a full messages channel.
	err := g.Wait()

	for i := 0; i < c.opts.Concurrency; i++ {
		c.stopWork <- struct{}{}
	}

	wg.Wait()

	close(c.Errors)

	return err
}

func (c *Consumer) poll(ctx context.Context, streams []string) error {
	for {
		if err := ctx.Err(); err != nil {
			return nil
		}

		results, err := c.xReadGroup(ctx, streams)
		if err != nil {
			return err
		}

		for _, result := range results {
			for _, x := range result.Messages {
				c.enqueue(result.Stream, &x, 0)
			}
		}
	}
}

func (c *Consumer) claim(ctx context.Context) error {
	if c.opts.ClaimInterval == 0 {
		return nil
	}

	for {
		for stream := range c.handlers {
			if err := c.claimStream(ctx, stream); err != nil {
				return err
			}
		}

		select {
		case <-time.After(c.opts.ClaimInterval):
			continue
		case <-ctx.Done():
			return nil
		}
	}
}

func (c *Consumer) claimStream(ctx context.Context, stream string) error {
	start, end := "-", "+"

	for {
		if err := ctx.Err(); err != nil {
			return nil
		}

		res, err := c.xPendingExt(ctx, stream, start, end)
		if err != nil {
			return err
		}

		// We're done; no more pending messages.
		if len(res) == 0 {
			break
		}

		toClaim := make([]string, 0, len(res))
		retryCnt := make(map[string]int64)

		for _, pending := range res {
			// Don't claim own messages
			if pending.Consumer == c.opts.Consumer {
				continue
			}

			if pending.Idle < c.opts.PendingTimeout {
				continue
			}

			toClaim = append(toClaim, pending.ID)
			retryCnt[pending.ID] = pending.RetryCount
		}

		claimed, err := c.xClaim(ctx, stream, toClaim)
		if err != nil {
			return err
		}

		for _, x := range claimed {
			c.enqueue(stream, &x, retryCnt[x.ID])
		}

		start, err = incrementMessageID(res[len(res)-1].ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Consumer) enqueue(stream string, x *redis.XMessage, rc int64) {
	c.messages <- &Message{
		ID:         x.ID,
		Stream:     stream,
		Values:     x.Values,
		retryCount: rc,
	}
}

func (c *Consumer) work(ctx context.Context) {
	for {
		select {
		case <-c.stopWork:
			return
		case msg := <-c.messages:
			if err := c.process(ctx, msg); err != nil {
				c.Errors <- err
			}
		}
	}
}

func (c *Consumer) process(ctx context.Context, msg *Message) (err error) {
	defer func() {
		if v := recover(); v != nil {
			if e, ok := v.(error); ok {
				err = &HandlerPanic{Err: e, Msg: msg}
				return
			}
			err = &HandlerPanic{Err: fmt.Errorf("%v", v)}
		}
	}()

	var handler Handler

	if msg.retryCount <= int64(c.opts.MaxRetry) || c.opts.MaxRetry < 0 {
		handler = c.handlers[msg.Stream]
	} else {
		handler = c.deadLetterHandlers[msg.Stream]
	}

	var r Result

	if handler != nil {
		handlerCtx, cancel := context.WithTimeout(ctx, c.opts.HandlerTimeout)
		defer cancel()
		r = handler.HandleEvent(handlerCtx, msg)
	}

	if r.Err != nil && r.Retry { // Failure
		// TODO should this be able to return fatal errors?
		if err := c.reschedule(ctx, msg, r); err != nil {
			return err
		}

		return &HandlerError{Err: r.Err, Msg: msg}
	}

	if err := c.xAck(ctx, msg.Stream, msg.ID); err != nil {
		return err
	}

	if r.Err != nil { // Discarded
		return &HandlerError{Err: r.Err, Msg: msg}
	}

	return nil // Success or no handler
}

func (c *Consumer) reschedule(ctx context.Context, msg *Message, res Result) error {
	d := res.Backoff
	if d == 0 {
		d = c.backoff.ForAttempt(float64(msg.retryCount))
	}

	select {
	case <-ctx.Done():
		return nil
	case <-time.After(d):
		// Continue
	}

	// Claim our own message to increment Redis' retry count
	claimed, err := c.xClaim(ctx, msg.Stream, []string{msg.ID})
	if err != nil {
		return err
	}

	if len(claimed) != 1 {
		return fmt.Errorf("failed to claim message %q from stream %s", msg.ID, msg.Stream)
	}

	pending, err := c.xPendingExt(ctx, msg.Stream, msg.ID, msg.ID)
	if err != nil {
		return err
	}

	if len(pending) != 1 {
		return fmt.Errorf("failed to read pending message %q from stream %q", msg.ID, msg.Stream)
	}

	for _, x := range claimed {
		c.enqueue(msg.Stream, &x, pending[0].RetryCount)
	}

	return nil
}

// ----- Redis functions -----
// These functions proxy through to the Redis client passed
// in the Options struct, but add retries with backoff.

func (c *Consumer) xGroupCreateMkStream(ctx context.Context, stream string) error {
	f := func() error {
		if result, err := c.opts.Redis.XGroupCreateMkStream(
			ctx,
			stream,
			c.opts.Group,
			"$", // Group consumes only new messages
		).Result(); err != nil {
			// It's fine if the consumer group already exists
			if err.Error() != errConsumerGroupAlreadyExists {
				return err
			}
		} else if result != "OK" {
			return fmt.Errorf("non-OK response from Redis: %s", result)
		}
		return nil
	}

	return c.xRetry(ctx, f)
}

func (c *Consumer) xReadGroup(ctx context.Context, streams []string) ([]redis.XStream, error) {
	var results []redis.XStream

	f := func() error {
		var err error
		results, err = c.opts.Redis.XReadGroup(ctx, &redis.XReadGroupArgs{
			Group:    c.opts.Group,
			Consumer: c.opts.Consumer,
			Streams:  streams,
			Count:    int64(c.opts.BufferSize - len(c.messages)),
			Block:    c.opts.ReadTimeout,
		}).Result()
		return err
	}

	return results, c.xRetry(ctx, f)
}

func (c *Consumer) xPendingExt(ctx context.Context, stream, start, end string) ([]redis.XPendingExt, error) {
	var results []redis.XPendingExt

	f := func() error {
		var err error
		results, err = c.opts.Redis.XPendingExt(ctx, &redis.XPendingExtArgs{
			Group:    c.opts.Group,
			Consumer: c.opts.Consumer,
			Stream:   stream,
			Start:    start,
			End:      end,
			Count:    int64(c.opts.BufferSize - len(c.messages)),
		}).Result()
		return err
	}

	return results, c.xRetry(ctx, f)
}

func (c *Consumer) xClaim(ctx context.Context, stream string, messages []string) ([]redis.XMessage, error) {
	if len(messages) == 0 {
		return nil, nil
	}

	var results []redis.XMessage

	f := func() error {
		var err error
		results, err = c.opts.Redis.XClaim(ctx, &redis.XClaimArgs{
			Stream:   stream,
			Group:    c.opts.Group,
			Consumer: c.opts.Consumer,
			MinIdle:  time.Minute,
			Messages: messages,
		}).Result()
		return err
	}

	return results, c.xRetry(ctx, f)
}

func (c *Consumer) xAck(ctx context.Context, stream, ID string) error {
	return c.xRetry(ctx, func() error {
		return c.opts.Redis.XAck(ctx, stream, c.opts.Group, ID).Err()
	})
}

func (c *Consumer) xRetry(ctx context.Context, f func() error) error {
	return xRetry(&xRetryArgs{
		ctx:      ctx,
		f:        f,
		errs:     c.Errors,
		b:        c.backoff,
		maxRetry: c.opts.NetworkRetry,
	})
}
