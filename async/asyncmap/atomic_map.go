// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asyncmap

import "sync"

// AtomicMap is a generic wrapper around the standard sync.Map.
//
// This map is optimized mostly for read operations.
// Use ShardMap if you need a map optimized for read-write operations.
type AtomicMap[K comparable, V any] interface {
	Map[K, V]
}

// NewAtomicMap returns a new atomic map backed by a sync.Map.
func NewAtomicMap[K comparable, V any]() AtomicMap[K, V] {
	return newAtomicMap[K, V]()
}

// internal

var _ AtomicMap[int, int] = (*atomicMap[int, int])(nil)

type atomicMap[K comparable, V any] struct {
	raw sync.Map
}

func newAtomicMap[K comparable, V any]() *atomicMap[K, V] {
	return &atomicMap[K, V]{}
}

// Len iterates the map, counts the number of keys, and returns the result.
func (m *atomicMap[K, V]) Len() int {
	var n int
	m.raw.Range(func(_, _ any) bool {
		n++
		return true
	})
	return n
}

// Clear deletes all items.
func (m *atomicMap[K, V]) Clear() {
	m.raw.Clear()
}

// Contains returns true if a key exists.
func (m *atomicMap[K, V]) Contains(key K) bool {
	_, ok := m.raw.Load(key)
	return ok
}

// Get returns a value by key, or false.
func (m *atomicMap[K, V]) Get(key K) (v V, _ bool) {
	val, ok := m.raw.Load(key)
	if !ok {
		return v, false
	}
	return val.(V), true
}

// GetOrSet returns a value by key, or sets a value if it does not exist.
func (m *atomicMap[K, V]) GetOrSet(key K, value V) (_ V, set bool) {
	val, loaded := m.raw.LoadOrStore(key, value)
	return val.(V), !loaded
}

// Delete deletes a value by key.
func (m *atomicMap[K, V]) Delete(key K) {
	m.raw.Delete(key)
}

// Pop deletes and returns a value by key, or false.
func (m *atomicMap[K, V]) Pop(key K) (v V, _ bool) {
	val, ok := m.raw.LoadAndDelete(key)
	if !ok {
		return v, false
	}
	return val.(V), true
}

// Set sets a value for a key.
func (m *atomicMap[K, V]) Set(key K, value V) {
	m.raw.Store(key, value)
}

// Swap swaps a key value and returns the previous value.
func (m *atomicMap[K, V]) Swap(key K, value V) (v V, _ bool) {
	val, ok := m.raw.Swap(key, value)
	if !ok {
		return v, false
	}
	return val.(V), true
}

// Range iterates over all key-value pairs.
// The iteration stops if the function returns false.
func (m *atomicMap[K, V]) Range(fn func(K, V) bool) {
	m.raw.Range(func(key, value any) bool {
		return fn(key.(K), value.(V))
	})
}
