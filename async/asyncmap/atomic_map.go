// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asyncmap

import (
	"sync"
	"sync/atomic"

	"github.com/basecomplextech/baselibrary/internal/hashing"
	"github.com/basecomplextech/baselibrary/pools"
)

// AtomicMap is a goroutine-safe hash map based on atomic operations.
// Readers are non-blocking, writers use a mutex per bucket, and a single resize mutex.
// The map is optimized for read/write operations and uses a pool for memory allocation.
//
// # Implementation
//
// The map uses a variant of reference counting with two reference counts
// described in "C++ Concurrency in Action" by Anthony Williams.
// Each entry has an external reference count and an internal reference count.
//
// Readers increment the external reference count when acquire an entry,
// and *decrement* the internal reference count when release an entry.
//
// Writes atomically swap the reference with a new one, and increment
// the internal reference count by the (external reference count - 1).
//
// When the internal reference count reaches zero, the entry is freed
// and returned to the pool.
//
// # Benchmarks
//
//	cpu: Apple M1 Pro
//	BenchmarkAtomicMap_Read-10                            	74086842	        14.26 ns/op	        70.15 mops	       0 B/op	       0 allocs/op
//	BenchmarkAtomicMap_Read_Parallel-10                   	61127808	        18.12 ns/op	        55.20 mops	       0 B/op	       0 allocs/op
//	BenchmarkAtomicMap_Write-10                           	27489913	        43.14 ns/op	        23.18 mops	       0 B/op	       0 allocs/op
//	BenchmarkAtomicMap_Write_Parallel-10                  	 7670359	       161.00 ns/op	         6.21 mops	       0 B/op	       0 allocs/op
//	BenchmarkAtomicMap_Read_Write_Parallel-10             	 7588120	       163.60 ns/op	         2.57 rmops	         6.111 wmops	       0 B/op	       0 allocs/op
//	BenchmarkAtomicMap_Read_Parallel_Write_Parallel-10    	 3187329	       360.30 ns/op	        43.66 rmops	         2.775 wmops	       0 B/op	       0 allocs/op
type AtomicMap[K comparable, V any] interface {
	Map[K, V]
}

// NewAtomicMap returns a new atomic map.
func NewAtomicMap[K comparable, V any]() AtomicMap[K, V] {
	return newAtomicMap[K, V](0)
}

// internal

const (
	// atomicMapThreshold is the load factor for resizing the map.
	atomicMapThreshold = 0.75

	// atomicMapMinSize is the minimum number of buckets.
	atomicMapMinSize = 16
)

var _ Map[int, int] = (*atomicMap[int, int])(nil)

type atomicMap[K comparable, V any] struct {
	pool pools.Pool[*atomicMapEntry[K, V]]

	wmu   sync.RWMutex                         // resize mutex
	state atomic.Pointer[atomicMapState[K, V]] // current state
}

func newAtomicMap[K comparable, V any](size int) *atomicMap[K, V] {
	pool := newAtomicMapPool[K, V]()

	num := int(float64(size) / atomicMapThreshold)
	num = max(num, atomicMapMinSize)

	s := newAtomicMapState(num, pool)
	m := &atomicMap[K, V]{pool: pool}
	m.state.Store(s)
	return m
}

func newAtomicMapPool[K comparable, V any]() pools.Pool[*atomicMapEntry[K, V]] {
	return pools.GetPool[*atomicMapEntry[K, V]](atomicMapEntryPools)
}

// Len returns the number of keys.
func (m *atomicMap[K, V]) Len() int {
	s := m.state.Load()
	return s.len()
}

// Clear deletes all items.
func (m *atomicMap[K, V]) Clear() {
	m.wmu.Lock()
	defer m.wmu.Unlock()

	s := m.state.Load()
	next := emptyAtomicMapState(s)

	m.state.Store(next)
}

// Contains returns true if a key exists.
func (m *atomicMap[K, V]) Contains(key K) bool {
	h := hashing.Hash(key)
	s := m.state.Load()
	return s.contains(h, key)
}

// Get returns a value by key, or false.
func (m *atomicMap[K, V]) Get(key K) (V, bool) {
	h := hashing.Hash(key)
	s := m.state.Load()
	return s.get(h, key)
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

// Delete deletes a key value, and returns the previous value.
func (m *atomicMap[K, V]) Delete(key K) (V, bool) {
	m.wmu.RLock()
	defer m.wmu.RUnlock()

	h := hashing.Hash(key)
	s := m.state.Load()
	return s.delete(h, key)
}

// Set sets a value for a key.
func (m *atomicMap[K, V]) Set(key K, value V) {
	resize := false
	m.set(key, value, &resize)

	if resize {
		m.resize()
	}
}

// SetAbsent sets a key value if absent, returns true if set.
func (m *atomicMap[K, V]) SetAbsent(key K, value V) bool {
	resize := false
	_, ok := m.getOrSet(key, value, &resize)

	if resize {
		m.resize()
	}
	return !ok
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

	h := hashing.Hash(key)
	s := m.state.Load()
	v, ok := s.getOrSet(h, key, value)

	n := s.len()
	*resize = n >= s.threshold
	return v, ok
}

func (m *atomicMap[K, V]) set(key K, value V, resize *bool) {
	m.wmu.RLock()
	defer m.wmu.RUnlock()

	h := hashing.Hash(key)
	s := m.state.Load()
	s.set(h, key, value)

	n := s.len()
	*resize = n >= s.threshold
}

func (m *atomicMap[K, V]) setAbsent(key K, value V, resize *bool) bool {
	m.wmu.RLock()
	defer m.wmu.RUnlock()

	h := hashing.Hash(key)
	s := m.state.Load()
	ok := s.setAbsent(h, key, value)

	n := s.len()
	*resize = n >= s.threshold
	return ok
}

func (m *atomicMap[K, V]) swap(key K, value V, resize *bool) (V, bool) {
	m.wmu.RLock()
	defer m.wmu.RUnlock()

	h := hashing.Hash(key)
	s := m.state.Load()
	v, ok := s.swap(h, key, value)

	n := s.len()
	*resize = n >= s.threshold
	return v, ok
}

// resize

func (m *atomicMap[K, V]) resize() {
	m.wmu.Lock()
	defer m.wmu.Unlock()

	s := m.state.Load()
	n := s.len()
	if n < s.threshold {
		return
	}

	// Double buckets
	size := len(s.buckets) * 2
	next := newAtomicMapState(size, s.pool)

	// Copy all items
	s.rangeLocked(func(k K, v V) bool {
		h := hashing.Hash(k)
		next.set(h, k, v)
		return true
	})

	// Replace state
	m.state.Store(next)
}

// pools

var atomicMapEntryPools = pools.NewPools()
