package async

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLockMap__should_lock_key(t *testing.T) {
	m := NewLockMap[int]()
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
	m := NewLockMap[int]()
	key := 123

	lock := m.Get(key)
	defer lock.Free()

	lock1, ok := m.locks[key]
	assert.True(t, ok)
	assert.Same(t, lock, lock1)
	assert.Equal(t, 1, lock1.refs)
}

func TestLockMap__should_retain_key_lock_when_already_locked(t *testing.T) {
	m := NewLockMap[int]()
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
	lock1, ok := m.locks[key]
	assert.True(t, ok)
	assert.Equal(t, 2, lock1.refs)
}

// Free

func TestKeyLock_Free__should_release_delete_key_lock(t *testing.T) {
	m := NewLockMap[int]()
	key := 123

	lock := m.Get(key)
	lock.Free()

	_, ok := m.locks[key]
	assert.False(t, ok)
}
