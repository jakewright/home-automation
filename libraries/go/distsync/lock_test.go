package distsync

import (
	"context"
	"testing"

	"gotest.tools/assert"
)

func TestLockSynchronous(t *testing.T) {
	DefaultLocksmith = NewLocalLocksmith()
	locker, err := Lock(context.Background(), "test")
	assert.NilError(t, err)
	locker.Unlock()

	locker, err = Lock(context.Background(), "test")
	assert.NilError(t, err)
	locker.Unlock()
}

func TestLockInterleaved(t *testing.T) {
	DefaultLocksmith = NewLocalLocksmith()
	locker, err := Lock(context.Background(), "test")
	assert.NilError(t, err)

	// Create a context that will immediately timeout
	ctx, cancel := context.WithTimeout(context.Background(), 0)
	defer cancel()

	locker2, err := Lock(ctx, "test")
	assert.Equal(t, nil, locker2)
	assert.ErrorContains(t, err, "failed to acquire lock in time")

	locker.Unlock()

	locker3, err := Lock(context.Background(), "test")
	assert.NilError(t, err)
	locker3.Unlock()
}
