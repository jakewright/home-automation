package dsync

import (
	"testing"

	"gotest.tools/assert"
)

func TestLockSynchronous(t *testing.T) {
	DefaultLocksmith = NewLocalLocksmith()
	locker, err := Lock("test")
	assert.NilError(t, err)
	locker.Unlock()

	locker, err = Lock("test")
	assert.NilError(t, err)
	locker.Unlock()
}

func TestLockInterleaved(t *testing.T) {
	DefaultLocksmith = NewLocalLocksmith()
	locker, err := Lock("test")
	assert.NilError(t, err)

	locker2, err := Lock("test")
	assert.Equal(t, nil, locker2)
	assert.ErrorContains(t, err, "Failed to acquire lock in time")

	locker.Unlock()

	locker3, err := Lock("test")
	assert.NilError(t, err)
	locker3.Unlock()
}
