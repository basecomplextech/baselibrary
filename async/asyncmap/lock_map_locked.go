// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asyncmap

// LockedLockMap is an interface to interact with a map locked in exclusive mode.
type LockedLockMap[K comparable] interface {
	// Contains returns true if the key is present.
	Contains(key K) bool

	// Range ranges over all keys.
	Range(fn func(key K) bool)

	// Free unlocks the map itself, internally it unlocks all buckets.
	Free()
}

// internal

var _ LockedLockMap[any] = (*lockedLockMap[any])(nil)

type lockedLockMap[K comparable] struct {
	m *lockMap[K]
}

func newLockedLockMap[K comparable](m *lockMap[K]) LockedLockMap[K] {
	return &lockedLockMap[K]{m: m}
}

// Contains returns true if the key is present.
func (m *lockedLockMap[K]) Contains(key K) bool {
	b := m.m.bucket(key)
	return b.containsLocked(key)
}

// Range ranges over all keys.
func (m *lockedLockMap[K]) Range(fn func(key K) bool) {
	for i := range m.m.buckets {
		b := &m.m.buckets[i]
		b.rangeLocked(fn)
	}
}

// Free unlocks the map itself, internally it unlocks all buckets.
func (m *lockedLockMap[K]) Free() {
	m.m.unlockMap()
	m.m = nil
}
