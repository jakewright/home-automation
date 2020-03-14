package dsync

import (
	"context"
	"fmt"
	"time"

	"github.com/jakewright/home-automation/libraries/go/errors"
	"github.com/jakewright/home-automation/libraries/go/slog"
)

// This package is a lie. It doesn't implement distributed locking at all.
// Distributed locking will be implemented if and when there are multiple
// instances of any services running.

const defaultTimeout = time.Second * 10

// Locker is a lock on a resource that can be locked and unlocked
type Locker interface {
	Lock() error
	Unlock()
}

// Locksmith can forge locks
type Locksmith interface {
	Forge(string) (Locker, error)
}

// DefaultLocksmith is a global instance of Locksmith
var DefaultLocksmith Locksmith

func mustGetDefaultLocksmith() Locksmith {
	if DefaultLocksmith == nil {
		slog.Panicf("dsync used before default locksmith set")
	}

	return DefaultLocksmith
}

// Lock will forge a lock for the resource and try to acquire the lock
func Lock(ctx context.Context, resource string, args ...interface{}) (Locker, error) {
	if ctx == nil {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), defaultTimeout)
		defer cancel()
	}

	for _, v := range args {
		resource = fmt.Sprintf("%s:%s", resource, v)
	}

	locker, err := mustGetDefaultLocksmith().Forge(resource)
	if err != nil {
		return nil, err
	}

	c := make(chan error)

	go func() {
		err := locker.Lock()

		select {
		case c <- err:
			// This is the normal case
		default:
			// There is no receiver which means the function
			// has already returned because the timeout was
			// reached. We don't need this lock anymore.
			locker.Unlock()
		}
	}()

	select {
	case err := <-c:
		if err != nil {
			return nil, errors.WithMessage(err, "failed to acquire lock")
		}
		return locker, nil
	case <-ctx.Done():
		return nil, errors.WithMessage(ctx.Err(), "failed to acquire lock in time")
	}
}
