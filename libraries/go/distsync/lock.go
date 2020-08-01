package distsync

import (
	"context"
	"fmt"
	"time"

	"github.com/jakewright/home-automation/libraries/go/slog"
)

const (
	defaultTimeout    = time.Second * 10
	defaultExpiration = time.Second * 60
)

// Locker is a lock on a resource that can be locked and unlocked
type Locker interface {
	Lock(ctx context.Context) error

	// Unlock releases the lock. Implementations of this can fail to unlock for
	// various reasons, so arguably this function should return an error. In
	// practice though, there's nothing useful you can do with the error and
	// the lock will expire at some point anyway.
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
	for _, v := range args {
		resource = fmt.Sprintf("%s:%s", resource, v)
	}

	locker, err := mustGetDefaultLocksmith().Forge(resource)
	if err != nil {
		return nil, err
	}

	return locker, locker.Lock(ctx)
}
