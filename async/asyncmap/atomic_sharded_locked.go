// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asyncmap

var _ LockedMap[any, any] = (*atomicShardedLocked[any, any])(nil)

type atomicShardedLocked[K comparable, V any] struct {
	m     *atomicShardedMap[K, V]
	freed bool
}

func newAtomicShardedLocked[K comparable, V any](m *atomicShardedMap[K, V]) LockedMap[K, V] {
	return &atomicShardedLocked[K, V]{m: m}
}

// Clear deletes all items.
func (m *atomicShardedLocked[K, V]) Clear() {
	if m.freed {
		panic("atomic map lock already unlocked")
	}

	for i := range m.m.shards {
		m.m.shards[i].clearLocked()
	}
}

// Range iterates over all key-value pairs.
// The iteration stops if the function returns false.
func (m *atomicShardedLocked[K, V]) Range(fn func(K, V) bool) {
	if m.freed {
		panic("atomic map lock already unlocked")
	}

	for i := range m.m.shards {
		ok := m.m.shards[i].range_(fn)
		if !ok {
			return
		}
	}
}

// Internal

// Free unlocks the locked map.
func (m *atomicShardedLocked[K, V]) Free() {
	if m.freed {
		panic("atomic map lock already unlocked")
	}
	m.freed = true

	for i := range m.m.shards {
		m.m.shards[i].wmu.Unlock()
	}
}
