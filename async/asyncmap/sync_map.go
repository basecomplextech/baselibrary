// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asyncmap

import "sync"

// SyncMap is a generic wrapper around the standard sync.Map.
//
// This map is optimized mostly for read operations.
// Writes operations are slower and allocate memory.
// Write performance degrades heavily if the keys are deleted from the map.
//
// Use [AtomicMap] or [AtomicShardedMap] if you need a map optimized for read-write operations.
//
// # Benchmarks
//
//	cpu: Apple M1 Pro
//	BenchmarkSyncMap_Read-10                            	38069961	        30.49 ns/op	        32.79 mops	       0 B/op	       0 allocs/op
//	BenchmarkSyncMap_Read_Parallel-10                   	274098883	         4.34 ns/op	       230.10 mops	       0 B/op	       0 allocs/op
//	BenchmarkSyncMap_Write-10                           	15118998	        79.18 ns/op	        12.63 mops	      28 B/op	       2 allocs/op
//	BenchmarkSyncMap_Write_Parallel-10                  	 4176376	       290.10 ns/op	         3.44 mops	      28 B/op	       2 allocs/op
//	BenchmarkSyncMap_Read_Write_Parallel-10             	31551691	        37.75 ns/op	        11.16 rmops	      26.49 wmops	      28 B/op	       2 allocs/op
//	BenchmarkSyncMap_Read_Parallel_Write_Parallel-10    	11116138	       126.10 ns/op	       174.20 rmops	       7.93 wmops	      28 B/op	       2 allocs/op
type SyncMap[K comparable, V any] interface {
	Map[K, V]
}

// NewSyncMap returns a new generic wrapper around a standard [sync.Map].
func NewSyncMap[K comparable, V any]() SyncMap[K, V] {
	return newSyncMap[K, V]()
}

// internal

var _ SyncMap[int, int] = (*syncMap[int, int])(nil)

type syncMap[K comparable, V any] struct {
	raw sync.Map
}

func newSyncMap[K comparable, V any]() *syncMap[K, V] {
	return &syncMap[K, V]{}
}

// Len iterates the map, counts the number of keys, and returns the result.
func (m *syncMap[K, V]) Len() int {
	var n int
	m.raw.Range(func(_, _ any) bool {
		n++
		return true
	})
	return n
}

// Clear deletes all items.
func (m *syncMap[K, V]) Clear() {
	m.raw.Clear()
}

// Contains returns true if a key exists.
func (m *syncMap[K, V]) Contains(key K) bool {
	_, ok := m.raw.Load(key)
	return ok
}

// Get returns a value by key, or false.
func (m *syncMap[K, V]) Get(key K) (v V, _ bool) {
	val, ok := m.raw.Load(key)
	if !ok {
		return v, false
	}
	return val.(V), true
}

// GetOrSet returns a value by key and true, or sets a value and false.
func (m *syncMap[K, V]) GetOrSet(key K, value V) (_ V, set bool) {
	val, ok := m.raw.LoadOrStore(key, value)
	return val.(V), ok
}

// Delete deletes a key value, and returns the previous value.
func (m *syncMap[K, V]) Delete(key K) (v V, _ bool) {
	val, ok := m.raw.LoadAndDelete(key)
	if !ok {
		return v, false
	}
	return val.(V), true
}

// LockMap is not supported.
func (m *syncMap[K, V]) LockMap() LockedMap[K, V] {
	panic("not supported")
}

// Set sets a value for a key.
func (m *syncMap[K, V]) Set(key K, value V) {
	m.raw.Store(key, value)
}

// SetAbsent sets a key value if absent, returns true if set.
func (m *syncMap[K, V]) SetAbsent(key K, value V) bool {
	_, loaded := m.raw.LoadOrStore(key, value)
	return !loaded
}

// Swap swaps a key value and returns the previous value.
func (m *syncMap[K, V]) Swap(key K, value V) (v V, _ bool) {
	val, ok := m.raw.Swap(key, value)
	if !ok {
		return v, false
	}
	return val.(V), true
}

// Range iterates over all key-value pairs.
// The iteration stops if the function returns false.
func (m *syncMap[K, V]) Range(fn func(K, V) bool) {
	m.raw.Range(func(key, value any) bool {
		return fn(key.(K), value.(V))
	})
}
