package async

import (
	"sync"
)

// LockMap holds locks for different keys.
//
// Usage:
//
//	m := NewLockMap[int]()
//
//	lock := m.Get(123)
//	defer lock.Free()
//
//	<-lock.Lock()
//	defer lock.Unlock()
type LockMap[K comparable] struct {
	mu    sync.Mutex
	locks map[K]*keyLock[K]
}

// KeyLock is a lock for a key in a lock map.
type KeyLock interface {
	// Lock returns a channel receiving from which locks the lock.
	Lock() <-chan struct{}

	// Unlock unlocks the lock.
	Unlock()

	// Free frees the acquired lock.
	Free()
}

// NewLockMap returns a new lock map.
func NewLockMap[K comparable]() *LockMap[K] {
	return &LockMap[K]{
		locks: make(map[K]*keyLock[K]),
	}
}

// Get returns a lock for the given key.
// The lock must be freed after use.
func (t *LockMap[K]) Get(key K) KeyLock {
	t.mu.Lock()
	defer t.mu.Unlock()

	lock, ok := t.locks[key]
	if ok {
		lock.refs++
	} else {
		lock = newKeyLock(t, key)
		t.locks[key] = lock
	}

	return lock
}

// internal

// release decrements the lock references, and deletes the lock if no references left.
func (t *LockMap[K]) release(lock *keyLock[K]) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if lock.refs <= 0 {
		panic("release of released lock")
	}

	lock.refs--
	if lock.refs > 0 {
		return
	}

	delete(t.locks, lock.key)
}

type keyLock[K comparable] struct {
	t    *LockMap[K]
	key  K
	lock Lock

	// guarded by map mutex
	refs int
}

func newKeyLock[K comparable](t *LockMap[K], key K) *keyLock[K] {
	return &keyLock[K]{
		t:    t,
		key:  key,
		lock: NewLock(),

		refs: 1,
	}
}

// Lock returns a channel receiving from which locks the lock.
func (l *keyLock[K]) Lock() <-chan struct{} {
	return l.lock
}

// Unlock unlocks the lock.
func (l *keyLock[K]) Unlock() {
	l.lock.Unlock()
}

// Free frees the acquired lock.
func (l *keyLock[K]) Free() {
	l.t.release(l)
}
