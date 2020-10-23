package queuer

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/jpillora/backoff"
)

// PublisherOptions contains options to configure the Publisher
type PublisherOptions struct {
	// StreamMaxLength sets the MAXLEN option when calling
	// XADD. This limits the size of the stream. Old entries
	// are automatically evicted when the specified length
	// is reached.
	StreamMaxLength int64

	// ApproximateMaxLength is an optimisation that allows
	// the stream to be capped more efficiently, as long as
	// an exact length is not required.
	ApproximateMaxLength bool

	// Redis is an instance of *redis.Client for use by
	// the client. This must be set before using the client.
	Redis *redis.Client

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

// Publisher is capable of publishing messages to Redis
// Streams. It should be created via NewPublisher().
type Publisher struct {
	// Errors is a channel over which non-fatal errors are
	// sent. This channel must have a listener otherwise
	// deadlock will arise.
	//
	// If NetworkRetry is set to 0, nothing will be sent
	// over this channel. It doesn't need a listener in
	// that case.
	Errors chan error

	opts    *PublisherOptions
	backoff *backoff.Backoff
}

// NewPublisher returns an initialised Publisher
func NewPublisher(opts *PublisherOptions) *Publisher {
	b := opts.Backoff
	if b == nil {
		b = defaultBackoff
	}

	return &Publisher{
		Errors:  make(chan error),
		opts:    opts,
		backoff: b,
	}
}

// Publish publishes the message to a Redis Stream
func (p *Publisher) Publish(ctx context.Context, m *Message) error {
	id, err := p.xAdd(ctx, m.Stream, m.ID, m.Values)
	if err != nil {
		return err
	}

	m.ID = id
	return nil
}

func (p *Publisher) xAdd(ctx context.Context, stream, id string, values map[string]interface{}) (string, error) {
	args := &redis.XAddArgs{
		Stream: stream,
		ID:     id,
		Values: values,
	}

	if p.opts.ApproximateMaxLength {
		args.MaxLenApprox = p.opts.StreamMaxLength
	} else {
		args.MaxLen = p.opts.StreamMaxLength
	}

	var newID string

	f := func() error {
		var err error
		newID, err = p.opts.Redis.XAdd(ctx, args).Result()
		return err
	}

	return newID, p.xRetry(ctx, f)
}

func (p *Publisher) xRetry(ctx context.Context, f func() error) error {
	return xRetry(&xRetryArgs{
		ctx:      ctx,
		f:        f,
		errs:     p.Errors,
		b:        p.backoff,
		maxRetry: p.opts.NetworkRetry,
	})
}
