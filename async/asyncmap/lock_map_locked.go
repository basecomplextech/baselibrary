// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asyncmap

// LockedMap is an interface to interact with a map locked in exclusive mode.
type LockedMap[K comparable] interface {
	// Contains returns true if the key is present.
	Contains(key K) bool

	// Range ranges over all keys.
	Range(fn func(key K) bool)

	// Free unlocks the map itself, internally it unlocks all buckets.
	Free()
}

// internal

var _ LockedMap[any] = (*lockedMap[any])(nil)

type lockedMap[K comparable] struct {
	m *lockMap[K]
}

func newLockedMap[K comparable](m *lockMap[K]) LockedMap[K] {
	return &lockedMap[K]{m: m}
}

// Contains returns true if the key is present.
func (m *lockedMap[K]) Contains(key K) bool {
	b := m.m.bucket(key)
	return b.containsLocked(key)
}

// Range ranges over all keys.
func (m *lockedMap[K]) Range(fn func(key K) bool) {
	for i := range m.m.buckets {
		b := &m.m.buckets[i]
		b.rangeLocked(fn)
	}
}

// Free unlocks the map itself, internally it unlocks all buckets.
func (m *lockedMap[K]) Free() {
	m.m.unlockMap()
	m.m = nil
}
