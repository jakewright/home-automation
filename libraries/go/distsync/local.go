package distsync

import (
	"context"
	"sync"

	"github.com/jakewright/home-automation/libraries/go/oops"
)

// LocalLocksmith implements process-scoped locking
type LocalLocksmith struct {
	locks sync.Map
}

// NewLocalLocksmith returns an initialised LocalLocksmith
func NewLocalLocksmith() *LocalLocksmith {
	return &LocalLocksmith{
		locks: sync.Map{},
	}
}

// Forge returns a Locker for the resource
func (l *LocalLocksmith) Forge(resource string) (Locker, error) {
	i, _ := l.locks.LoadOrStore(resource, &sync.Mutex{})
	mu := i.(*sync.Mutex)

	return &mutexWrapper{mu}, nil
}

type mutexWrapper struct {
	mu *sync.Mutex
}

// Lock acquires the lock
func (mw *mutexWrapper) Lock(ctx context.Context) error {
	if mw == nil {
		return oops.InternalService("tried to lock a nil locker")
	}

	c := make(chan struct{})

	go func() {
		mw.mu.Lock()

		select {
		case c <- struct{}{}:
			// This is the normal case
		default:
			// There is no receiver which means the function
			// has already returned because the timeout was
			// reached. We don't need this lock anymore.
			mw.Unlock()
		}
	}()

	select {
	case <-c:
		return nil // Lock acquired
	case <-ctx.Done():
		return oops.WithMessage(ctx.Err(), "failed to acquire lock in time")
	}
}

// Unlock releases the lock
func (mw *mutexWrapper) Unlock() {
	if mw == nil {
		return // probably ok ¯\_(ツ)_/¯
	}
	mw.mu.Unlock()
}
