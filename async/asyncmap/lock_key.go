// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asyncmap

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

// internal

var _ KeyLock = &keyLock[any]{}

type keyLock[K comparable] struct {
	item *lockItem[K]
}

// Lock returns a channel receiving from which locks the key.
func (l *keyLock[K]) Lock() <-chan struct{} {
	return l.item.lock
}

// Unlock unlocks the key lock.
func (l *keyLock[K]) Unlock() {
	l.item.unlock()
}

// Free frees the acquired key.
func (l *keyLock[K]) Free() {
	if l.item == nil {
		panic("free of freed key lock")
	}

	item := l.item
	l.item = nil
	item.free()
}

// locked

var _ LockedKey = &lockedKey[any]{}

type lockedKey[K comparable] struct {
	item *lockItem[K]
}

func (l *lockedKey[K]) Free() {
	if l.item == nil {
		panic("free of freed locked key")
	}

	item := l.item
	l.item = nil

	item.unlock()
	item.free()
}
