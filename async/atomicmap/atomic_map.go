// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package atomics

import (
	"sync"
	"sync/atomic"

	"github.com/basecomplextech/baselibrary/async/asyncmap"
	"github.com/basecomplextech/baselibrary/pools"
)

const (
	// atomicMapThreshold is the load factor for resizing the map.
	atomicMapThreshold = 0.75

	// atomicMapMinSize is the minimum number of buckets.
	atomicMapMinSize = 16
)

var _ asyncmap.Map[int, int] = (*atomicMap[int, int])(nil)

type atomicMap[K comparable, V any] struct {
	pool pools.Pool[*atomicEntry[K, V]]

	wmu   sync.RWMutex                      // resize mutex
	state atomic.Pointer[atomicState[K, V]] // current state
}

func newAtomicMap[K comparable, V any](size int) *atomicMap[K, V] {
	pool := newAtomicPool[K, V]()

	num := int(float64(size) / atomicMapThreshold)
	num = max(num, atomicMapMinSize)

	s := newAtomicState(num, pool)
	m := &atomicMap[K, V]{pool: pool}
	m.state.Store(s)
	return m
}

func newAtomicPool[K comparable, V any]() pools.Pool[*atomicEntry[K, V]] {
	return pools.NewPoolFunc(
		func() *atomicEntry[K, V] {
			return &atomicEntry[K, V]{}
		},
	)
}

// Len returns the number of keys.
func (m *atomicMap[K, V]) Len() int {
	s := m.state.Load()
	return int(s.count)
}

// Clear deletes all items.
func (m *atomicMap[K, V]) Clear() {
	m.wmu.Lock()
	defer m.wmu.Unlock()

	s := m.state.Load()
	next := clearAtomicState(s)

	m.state.Store(next)
}

// Contains returns true if a key exists.
func (m *atomicMap[K, V]) Contains(key K) bool {
	s := m.state.Load()
	return s.contains(key)
}

// Get returns a value by key, or false.
func (m *atomicMap[K, V]) Get(key K) (V, bool) {
	s := m.state.Load()
	return s.get(key)
}

// GetOrSet returns a value by key and true, or sets a value and false.
func (m *atomicMap[K, V]) GetOrSet(key K, value V) (V, bool) {
	resize := false
	v, ok := m.getOrSet(key, value, &resize)

	if resize {
		m.resize()
	}
	return v, ok
}

// Delete deletes a value by key.
func (m *atomicMap[K, V]) Delete(key K) {
	m.wmu.RLock()
	defer m.wmu.RUnlock()

	s := m.state.Load()
	s.delete(key)
}

// Pop deletes and returns a value by key, or false.
func (m *atomicMap[K, V]) Pop(key K) (V, bool) {
	m.wmu.RLock()
	defer m.wmu.RUnlock()

	s := m.state.Load()
	return s.pop(key)
}

// Set sets a value for a key.
func (m *atomicMap[K, V]) Set(key K, value V) {
	resize := false
	m.set(key, value, &resize)

	if resize {
		m.resize()
	}
}

// Swap swaps a key value and returns the previous value.
func (m *atomicMap[K, V]) Swap(key K, value V) (V, bool) {
	resize := false
	v, ok := m.swap(key, value, &resize)

	if resize {
		m.resize()
	}
	return v, ok
}

// Range iterates over all key-value pairs.
// The iteration stops if the function returns false.
func (m *atomicMap[K, V]) Range(fn func(K, V) bool) {
	s := m.state.Load()
	s.range_(fn)
}

// private

func (m *atomicMap[K, V]) getOrSet(key K, value V, resize *bool) (V, bool) {
	m.wmu.RLock()
	defer m.wmu.RUnlock()

	s := m.state.Load()
	v, ok := s.getOrSet(key, value)

	*resize = s.count >= int64(s.threshold)
	return v, ok
}

func (m *atomicMap[K, V]) set(key K, value V, resize *bool) {
	m.wmu.RLock()
	defer m.wmu.RUnlock()

	s := m.state.Load()
	s.set(key, value)

	*resize = s.count >= int64(s.threshold)
}

func (m *atomicMap[K, V]) swap(key K, value V, resize *bool) (V, bool) {
	m.wmu.RLock()
	defer m.wmu.RUnlock()

	s := m.state.Load()
	v, ok := s.swap(key, value)

	*resize = s.count >= int64(s.threshold)
	return v, ok
}

// resize

func (m *atomicMap[K, V]) resize() {
	m.wmu.Lock()
	defer m.wmu.Unlock()

	s := m.state.Load()
	if s.count < int64(s.threshold) {
		return
	}

	next := resizeAtomicState(s)
	m.state.Store(next)
}
