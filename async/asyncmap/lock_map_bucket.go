// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asyncmap

import (
	"sync"
)

type lockMapBucket[K comparable] struct {
	m *lockMap[K]

	mu    sync.Mutex
	entry lockMapEntry[K]
}

func (b *lockMapBucket[K]) init(m *lockMap[K]) {
	b.m = m
}

// get returns an item by a key, increments its refs, or inserts a new one.
func (b *lockMapBucket[K]) get(key K) *lockMapItem[K] {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Get existing, increment refs
	if m, ok := b.entry.get(key); ok {
		m.refs++
		return m
	}

	// Add new item
	m := newLockMapItem[K](b, key)
	b.entry.set(m)
	return m
}

// getNoRetain returns an item without incrementing refs, used in tests.
func (b *lockMapBucket[K]) getNoRetain(key K) (*lockMapItem[K], bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.entry.get(key)
}

// contains

// contains returns true if a key exists.
func (b *lockMapBucket[K]) contains(key K) bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	return b.entry.contains(key)
}

// release

// release decrements item refs, and deletes the item when refs reach 0.
func (b *lockMapBucket[K]) release(m *lockMapItem[K]) (deleted bool) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Decrement refs
	if m.refs <= 0 {
		panic("free of freed key lock")
	}
	m.refs--
	if m.refs > 0 {
		return false
	}

	// Delete item
	b.entry.delete(m)
	return true
}

// locked

// containsLocked returns true if a key exists, must be called with lock held.
func (b *lockMapBucket[K]) containsLocked(key K) bool {
	return b.entry.contains(key)
}

// rangeLocked calls a function for each key, must be called with lock held.
func (b *lockMapBucket[K]) rangeLocked(fn func(key K) bool) {
	b.entry.range_(fn)
}
