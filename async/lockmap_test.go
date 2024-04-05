package async

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLockMap__should_lock_key(t *testing.T) {
	m := newLockMap[int]()
	key := 123
	value := 0

	lock := m.Get(key)
	defer lock.Free()

	<-lock.Lock()

	done := make(chan struct{})
	go func() {
		defer close(done)

		lock := m.Get(key)
		defer lock.Free()

		<-lock.Lock()
		defer lock.Unlock()

		value = 3
	}()

	time.Sleep(10 * time.Millisecond)
	value = 2
	lock.Unlock()

	<-done
	assert.Equal(t, 3, value)
}

func TestLockMap__should_retain_key_lock(t *testing.T) {
	m := newLockMap[int]()
	key := 123

	lock := m.Get(key)
	defer lock.Free()

	lock1_, ok := m.locks.Load(key)
	require.True(t, ok)

	lock1 := lock1_.(*keyLock[int])
	lock1.mu.Lock() // pass race detector
	defer lock1.mu.Unlock()

	assert.Same(t, lock, lock1)
	assert.Equal(t, int32(1), lock1.refs)
}

func TestLockMap__should_retain_key_lock_when_already_locked(t *testing.T) {
	m := newLockMap[int]()
	key := 123

	lock := m.Get(key)
	defer lock.Free()

	<-lock.Lock()
	defer lock.Unlock()

	go func() {
		lock := m.Get(key)
		defer lock.Free()

		<-lock.Lock()
		lock.Unlock()
	}()

	time.Sleep(10 * time.Millisecond)
	lock1_, ok := m.locks.Load(key)
	require.True(t, ok)

	lock1 := lock1_.(*keyLock[int])
	lock1.mu.Lock() // pass race detector
	defer lock1.mu.Unlock()

	assert.Equal(t, int32(2), lock1.refs)
}

// Lock

func TestLockMap_Lock__should_acquire_locked_key(t *testing.T) {
	m := newLockMap[int]()
	key := 123
	ctx := NoContext()

	lock, st := m.Lock(ctx, key)
	if !st.OK() {
		t.Fatal(st)
	}
	lock.Free()

	_, ok := m.locks.Load(key)
	assert.False(t, ok)
}

// Free

func TestKeyLock_Free__should_release_delete_key_lock(t *testing.T) {
	m := newLockMap[int]()
	key := 123

	lock := m.Get(key)
	lock.Free()

	_, ok := m.locks.Load(key)
	assert.False(t, ok)
}
