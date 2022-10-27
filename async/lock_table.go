package async

import (
	"sync"

	"github.com/complex1tech/baselibrary/status"
)

// LockTable holds locks for different keys.
type LockTable[K comparable] struct {
	mu    sync.Mutex
	locks map[K]*keyLock[K]
}

// NewLockTable returns a new lock table.
func NewLockTable[K comparable]() *LockTable[K] {
	return &LockTable[K]{
		locks: make(map[K]*keyLock[K]),
	}
}

// Lock locks a key and returns the key lock.
func (t *LockTable[K]) Lock(cancel <-chan struct{}, key K) (KeyLock, status.Status) {
	// acquire key lock
	lock, st := t.acquire(key)
	if !st.OK() {
		return nil, st
	}
	defer t.release(lock)

	// lock key lock
	select {
	case <-lock.lock:
	case <-cancel:
		return nil, status.Cancelled
	}

	t.retain(lock)
	return lock, status.OK
}

// acquire returns a lock for the given key, increments its references.
func (t *LockTable[K]) acquire(key K) (*keyLock[K], status.Status) {
	t.mu.Lock()
	defer t.mu.Unlock()

	lock, ok := t.locks[key]
	if ok {
		lock.refs++
	} else {
		lock = newKeyLock(t, key)
		t.locks[key] = lock
	}

	return lock, status.OK
}

// retain increments the lock references, panics if already released.
func (t *LockTable[K]) retain(lock *keyLock[K]) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if lock.refs <= 0 {
		panic("retain of released lock")
	}

	lock.refs++
}

// release decrements the lock references, and deletes the lock if no references left.
func (t *LockTable[K]) release(lock *keyLock[K]) {
	t.mu.Lock()
	defer t.mu.Unlock()

	lock.refs--
	if lock.refs > 0 {
		return
	}

	delete(t.locks, lock.key)
}

// Lock

// KeyLock is a lock for a key in a lock table.
type KeyLock interface {
	// Unlock unlocks the lock.
	Unlock()
}

type keyLock[K comparable] struct {
	t    *LockTable[K]
	key  K
	lock Lock

	// guarded by table mutex
	refs int
}

func newKeyLock[K comparable](t *LockTable[K], key K) *keyLock[K] {
	return &keyLock[K]{
		t:    t,
		key:  key,
		lock: NewLock(),

		refs: 1,
	}
}

// Unlock unlocks the lock.
func (l *keyLock[K]) Unlock() {
	l.lock.Unlock()
	l.t.release(l)
}
