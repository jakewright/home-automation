package queuer

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/jpillora/backoff"
)

func init() {
	// Used by the backoff library
	rand.Seed(time.Now().UnixNano())
}

var defaultBackoff = &backoff.Backoff{
	Min:    100 * time.Millisecond,
	Max:    10 * time.Second,
	Factor: 2,
	Jitter: true,
}

// incrementMessageID takes in a message ID (e.g. 1564886140363-0) and
// increments the index section (e.g. 1564886140363-1). This is the next valid
// ID value, and it can be used for paging through messages.
// https://github.com/robinjoseph08/redisqueue/blob/195b427f6d5d99c0986d9b219811155ecfc01ee8/redis.go#L55-L66
func incrementMessageID(id string) (string, error) {
	parts := strings.Split(id, "-")
	index := parts[1]
	parsed, err := strconv.ParseInt(index, 10, 64)
	if err != nil {
		return "", fmt.Errorf("error parsing message ID %q: %w", id, err)
	}
	return fmt.Sprintf("%s-%d", parts[0], parsed+1), nil
}

type xRetryArgs struct {
	ctx      context.Context
	f        func() error
	errs     chan<- error
	b        *backoff.Backoff
	maxRetry int
}

// xRetry will call f(), retrying in the case of network
// errors with backoff b a maximum of maxRetry times or
// until the context is cancelled.
func xRetry(args *xRetryArgs) error {
	// The select statement below is not guaranteed to return
	// a non-nil context error in the first iteration because
	// go chooses a case at random if they're both "ready".
	if args.ctx.Err() != nil {
		return args.ctx.Err()
	}

	var sleep time.Duration

	for i := 0; ; i++ {
		select {
		case <-args.ctx.Done():
			return args.ctx.Err()
		case <-time.After(sleep):
			// Continue
		}

		err := args.f()
		if err == nil {
			return nil
		}

		// Was this a network error?
		if err, ok := err.(net.Error); ok {
			// Ignore temporary errors
			if err.Timeout() || err.Temporary() {
				// Have we reached the retry limit already?
				if i == args.maxRetry {
					return &NetworkError{Err: err}
				}

				sleep = args.b.ForAttempt(float64(i))
				args.errs <- &NetworkError{Err: err, Retrying: true, Backoff: sleep}
				continue
			}

			// Anything else is a fatal error
			return &NetworkError{Err: err}
		}

		return &RedisError{Err: err}
	}
}
