package async

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLockTable_Lock__should_lock_key(t *testing.T) {
	tbl := NewLockTable[int]()
	key := 123
	value := 0

	lock, st := tbl.Lock(nil, key)
	if !st.OK() {
		t.Fatal(st)
	}

	done := make(chan struct{})
	go func() {
		defer close(done)

		lock, st := tbl.Lock(nil, key)
		if !st.OK() {
			t.Fatal(st)
		}
		defer lock.Unlock()

		value = 3
	}()

	time.Sleep(10 * time.Millisecond)
	value = 2
	lock.Unlock()

	<-done
	assert.Equal(t, 3, value)
}

func TestLockTable_Lock__should_retain_key_lock(t *testing.T) {
	tbl := NewLockTable[int]()
	key := 123

	lock, st := tbl.Lock(nil, key)
	if !st.OK() {
		t.Fatal(st)
	}
	defer lock.Unlock()

	lock1, ok := tbl.locks[key]
	assert.True(t, ok)
	assert.Same(t, lock, lock1)
	assert.Equal(t, 1, lock1.refs)
}

func TestLockTable_Lock__should_retain_key_lock_when_already_locked(t *testing.T) {
	tbl := NewLockTable[int]()
	key := 123

	lock, st := tbl.Lock(nil, key)
	if !st.OK() {
		t.Fatal(st)
	}
	defer lock.Unlock()

	go func() {
		lock, st := tbl.Lock(nil, key)
		if !st.OK() {
			t.Fatal(st)
		}
		lock.Unlock()
	}()

	time.Sleep(10 * time.Millisecond)
	lock1, ok := tbl.locks[key]
	assert.True(t, ok)
	assert.Equal(t, 2, lock1.refs)
}

// Unlock

func TestLockTable_Unlock__should_release_delete_key_lock(t *testing.T) {
	tbl := NewLockTable[int]()
	key := 123

	lock, st := tbl.Lock(nil, key)
	if !st.OK() {
		t.Fatal(st)
	}
	lock.Unlock()

	_, ok := tbl.locks[key]
	assert.False(t, ok)
}

func TestLockTable_Unlock__should_release_key_lock(t *testing.T) {
	tbl := NewLockTable[int]()
	key := 123

	lock, st := tbl.Lock(nil, key)
	if !st.OK() {
		t.Fatal(st)
	}

	proceed := make(chan struct{})
	go func() {
		lock, st := tbl.Lock(nil, key)
		if !st.OK() {
			t.Fatal(st)
		}

		<-proceed
		lock.Unlock()
	}()
	defer close(proceed)

	time.Sleep(10 * time.Millisecond)
	lock1, ok := tbl.locks[key]
	assert.True(t, ok)
	assert.Equal(t, 2, lock1.refs)

	lock.Unlock()
	assert.Equal(t, 1, lock1.refs)
}
