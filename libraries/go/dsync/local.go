package dsync

import (
	"sync"
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

func (mw *mutexWrapper) Lock() error {
	mw.mu.Lock()
	return nil
}

func (mw *mutexWrapper) Unlock() {
	mw.mu.Unlock()
}
