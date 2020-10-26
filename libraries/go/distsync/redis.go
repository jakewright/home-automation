package distsync

import (
	"context"
	"fmt"
	"time"

	"github.com/danielchatfield/go-randutils"
	"github.com/go-redis/redis/v8"

	"github.com/jakewright/home-automation/libraries/go/oops"
	"github.com/jakewright/home-automation/libraries/go/slog"
)

var luaReleaseLock = redis.NewScript(`
if redis.call("get", KEYS[1]) == ARGV[1] then
    return redis.call("del", KEYS[1])
else
    return 0
end
`)

// redisLock is a lock backed by a single Redis node
type redisLock struct {
	key        string
	value      string
	timeout    time.Duration
	expiration time.Duration
	client     *redis.Client
}

var _ Locker = (*redisLock)(nil)

// Lock obtains the lock
func (l *redisLock) Lock(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, l.timeout)
	defer cancel()

	for {
		if err := ctx.Err(); err != nil {
			return oops.WithMessage(err, "failed to acquire lock in time")
		}

		set, err := l.client.SetNX(ctx, l.key, l.value, l.expiration).Result()
		if err != nil {
			// We don't know whether we have the lock in this case, but even if we
			// do, it'll expire at some point anyway so there's no need to try to
			// unlock. Doing so could also return an error making this not trivial.
			return oops.WithMessage(err, "failed to acquire lock")
		}

		if set {
			return nil
		}
	}
}

// Unlock releases the lock
func (l *redisLock) Unlock() {
	rsp, err := luaReleaseLock.Run(context.TODO(), l.client, []string{l.key}, l.value).Result()
	if err != nil && err != redis.Nil {
		slog.Errorf("Failed to release lock on resource %q", l.key)
	}

	// The response is the number of keys that were deleted. If the lock has
	// already been unlocked, this will be zero. To match the behaviour of
	// the sync package, panic in this case.
	switch rsp {
	case 0:
		panic("nothing to unlock")
	case 1: // ok
	default:
		panic("unexpected number of records deleted")
	}
}

// RedisLocksmith is a Locksmith backed by Redis. It is important that there is
// only a single Redis node. This algorithm will not work with a Redis cluster.
type RedisLocksmith struct {
	// ServiceName is used to scope locks to the current service only
	ServiceName string

	// Client is a Redis client
	Client *redis.Client

	// Timeout is the length of time to wait when trying to
	// acquire a lock before giving up and returning an error
	Timeout time.Duration

	// Expiration is the TTL to set on locks
	Expiration time.Duration
}

var _ Locksmith = (*RedisLocksmith)(nil)

// Forge returns a Locker that can be locked and unlocked
func (l *RedisLocksmith) Forge(resource string) (Locker, error) {
	errParams := map[string]string{
		"resource": resource,
	}

	val, err := randutils.String(16)
	if err != nil {
		return nil, oops.WithMessage(err, "failed to generate random value", errParams)
	}

	timeout := l.Timeout
	if timeout == 0 {
		timeout = defaultTimeout
	}

	expiration := l.Expiration
	if expiration == 0 {
		expiration = defaultExpiration
	}

	return &redisLock{
		key:        fmt.Sprintf("%s:%s", l.ServiceName, resource),
		value:      val,
		timeout:    timeout,
		expiration: expiration,
		client:     l.Client,
	}, nil
}
