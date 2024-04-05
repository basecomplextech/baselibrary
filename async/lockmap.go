package async

import (
	"sync"

	"github.com/basecomplextech/baselibrary/status"
)

// LockMap holds locks for different keys.
type LockMap[K comparable] interface {
	// Get returns a key lock, the lock must be freed after use.
	//
	// Usage:
	//
	//	m := NewLockMap[int]()
	//
	//	lock := m.Get(123)
	//	defer lock.Free()
	//
	//	select {
	//	case <-lock.Lock():
	//	case <-time.After(time.Second):
	//		return status.Timeout
	//	case <-ctx.Wait():
	//		return ctx.Status()
	//	}
	//	defer lock.Unlock()
	Get(key K) KeyLock

	// Lock returns a locked key, the key must be freed after use.
	//
	// Usage:
	//
	//	m := NewLockMap[int]()
	//
	//	lock, st := m.Lock(ctx, 123)
	//	if !st.OK() {
	//		return st
	//	}
	//	defer lock.Free()
	Lock(ctx Context, key K) (LockedKey, status.Status)
}

// KeyLock is a single lock for a key, the lock must be freed after use.
type KeyLock interface {
	// Lock returns a channel receiving from which locks the key.
	Lock() <-chan struct{}

	// Unlock unlocks the key lock.
	Unlock()

	// Free frees the acquired key.
	Free()
}

// LockedKey is a locked key which is unlocked when freed.
type LockedKey interface {
	// Free unlocks and freed the key.
	Free()
}

// NewLockMap returns a new lock map.
func NewLockMap[K comparable]() LockMap[K] {
	return newLockMap[K]()
}

// internal

var _ LockMap[any] = &lockMap[any]{}

type lockMap[K comparable] struct {
	locks sync.Map
}

func newLockMap[K comparable]() *lockMap[K] {
	return &lockMap[K]{}
}

// Get returns a key key, the lock must be freed after use.
func (m *lockMap[K]) Get(key K) KeyLock {
	for {
		if l, ok := m.load(key); ok {
			return l
		}
		if l, ok := m.store(key); ok {
			return l
		}
	}
}

// Lock returns a locked key, the key must be freed after use.
func (m *lockMap[K]) Lock(ctx Context, key K) (LockedKey, status.Status) {
	// Get key lock
	l := m.Get(key)

	// Free if not locked
	ok := false
	defer func() {
		if !ok {
			l.Free()
		}
	}()

	// Try to lock
	select {
	case <-l.Lock():
	case <-ctx.Wait():
		return nil, ctx.Status()
	}

	// Return locked key
	k := newLockedKey[K](l.(*keyLock[K]))
	ok = true
	return k, status.OK
}

// private

func (m *lockMap[K]) load(key K) (*keyLock[K], bool) {
	// Try to load lock
	lock, ok := m.locks.Load(key)
	if !ok {
		return nil, false
	}

	// Try to increment refs
	l := lock.(*keyLock[K])
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.refs == 0 {
		return nil, false
	}

	l.refs++
	return l, true
}

func (m *lockMap[K]) store(key K) (*keyLock[K], bool) {
	// Make lock with 1 ref
	l := newKeyLock(m, key)
	l.mu.Lock()
	defer l.mu.Unlock()

	// Try to store it
	_, loaded := m.locks.LoadOrStore(key, l)
	if loaded {
		return nil, false
	}
	return l, true
}

func (m *lockMap[K]) free(l *keyLock[K]) {
	// Lock mutex
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.refs <= 0 {
		panic("free of freed key lock")
	}

	// Decrement refs
	l.refs--
	if l.refs > 0 {
		return
	}

	// Delete lock with 0 refs, while holding the mutex.
	// Any other goroutine, which already loaded it, will read 0 refs and skip it.
	m.locks.Delete(l.key)
}

// key

var _ KeyLock = &keyLock[any]{}

type keyLock[K comparable] struct {
	m    *lockMap[K]
	key  K
	lock chan struct{}

	mu   sync.Mutex
	refs int32
}

func newKeyLock[K comparable](m *lockMap[K], key K) *keyLock[K] {
	l := &keyLock[K]{
		m:    m,
		key:  key,
		lock: make(chan struct{}, 1),
		refs: 1,
	}
	l.lock <- struct{}{}
	return l
}

// Lock returns a channel receiving from which locks the key.
func (l *keyLock[K]) Lock() <-chan struct{} {
	return l.lock
}

// Unlock unlocks the key lock.
func (l *keyLock[K]) Unlock() {
	select {
	case l.lock <- struct{}{}:
	default:
		panic("unlock of unlocked key lock")
	}
}

// Free frees the acquired key.
func (l *keyLock[K]) Free() {
	l.m.free(l)
}

// lockedKey

var _ LockedKey = &lockedKey[any]{}

type lockedKey[K comparable] struct {
	lock *keyLock[K]
}

func newLockedKey[K comparable](lock *keyLock[K]) *lockedKey[K] {
	return &lockedKey[K]{lock: lock}
}

func (l *lockedKey[K]) Free() {
	defer l.lock.Free()

	l.lock.Unlock()
}
