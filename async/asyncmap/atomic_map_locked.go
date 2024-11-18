// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asyncmap

var _ LockedMap[any, any] = (*lockedAtomicMap[any, any])(nil)

type lockedAtomicMap[K comparable, V any] struct {
	m     *atomicMap[K, V]
	freed bool
}

func newLockedAtomicMap[K comparable, V any](m *atomicMap[K, V]) LockedMap[K, V] {
	return &lockedAtomicMap[K, V]{m: m}
}

// Clear deletes all items.
func (m *lockedAtomicMap[K, V]) Clear() {
	if m.freed {
		panic("atomic map lock already unlocked")
	}

	m.m.clearLocked()
}

// Range iterates over all key-value pairs.
// The iteration stops if the function returns false.
func (m *lockedAtomicMap[K, V]) Range(fn func(K, V) bool) {
	if m.freed {
		panic("atomic map lock already unlocked")
	}

	m.m.Range(fn)
}

// Internal

// Free unlocks the locked map.
func (m *lockedAtomicMap[K, V]) Free() {
	if m.freed {
		panic("atomic map lock already unlocked")
	}

	m.freed = true
	m.m.wmu.Unlock()
}
