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

// internal

var _ KeyLock = &lockMapKeyLock[any]{}

type lockMapKeyLock[K comparable] struct {
	item *lockMapItem[K]
}

func newLockMapKeyLock[K comparable](item *lockMapItem[K]) *lockMapKeyLock[K] {
	return &lockMapKeyLock[K]{item: item}
}

// Lock returns a channel receiving from which locks the key.
func (l *lockMapKeyLock[K]) Lock() <-chan struct{} {
	return l.item.lock
}

// Unlock unlocks the key lock.
func (l *lockMapKeyLock[K]) Unlock() {
	l.item.unlock()
}

// Free frees the acquired key.
func (l *lockMapKeyLock[K]) Free() {
	if l.item == nil {
		panic("free of freed key lock")
	}

	m := l.item
	l.item = nil

	m.release()
}
