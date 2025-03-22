// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asyncmap

import (
	"runtime"

	"github.com/basecomplextech/baselibrary/hashing"
)

// ShardedMap is a map which uses multiple hash maps each guarded by a separate mutex.
//
// This map is optimized for read-write operations.
// Use [SyncMap] if you need a map optimized mostly for read operations.
// Use [AtomicMap] or [AtomicShardedMap] if you need a map optimized for read/write
// operations and non-blocking reads.
//
// # Benchmarks
//
//	cpu: Apple M1 Pro
//	BenchmarkShardedMap_Read-10                            	28416379	        37.92 ns/op	        26.37 mops	       0 B/op	       0 allocs/op
//	BenchmarkShardedMap_Read_Parallel-10                   	45233161	        25.85 ns/op	        38.68 mops	       0 B/op	       0 allocs/op
//	BenchmarkShardedMap_Write-10                           	28739761	        41.63 ns/op	        24.02 mops	       0 B/op	       0 allocs/op
//	BenchmarkShardedMap_Write_Parallel-10                  	11624372	       100.70 ns/op	         9.92 mops	       0 B/op	       0 allocs/op
//	BenchmarkShardedMap_Read_Write_Parallel-10             	 8651703	       139.70 ns/op	         1.31 rmops	       7.15 wmops	       0 B/op	       0 allocs/op
//	BenchmarkShardedMap_Read_Parallel_Write_Parallel-10    	 1000000	      1162.00 ns/op	        22.56 rmops	       0.86 wmops	       0 B/op	       0 allocs/op
type ShardedMap[K comparable, V any] interface {
	Map[K, V]
}

// NewShardedMap returns a new sharded map.
func NewShardedMap[K comparable, V any]() Map[K, V] {
	return newShardedMap[K, V]()
}

// internal

var _ Map[int, int] = &shardedMap[int, int]{}

type shardedMap[K comparable, V any] struct {
	shards []shardedMapShard[K, V]
}

func newShardedMap[K comparable, V any]() *shardedMap[K, V] {
	cpus := runtime.NumCPU()

	return &shardedMap[K, V]{
		shards: make([]shardedMapShard[K, V], cpus),
	}
}

// Len returns the number of keys.
func (m *shardedMap[K, V]) Len() int {
	var n int

	for i := range m.shards {
		s := &m.shards[i]
		n += s.len()
	}

	return n
}

// Clear deletes all items.
func (m *shardedMap[K, V]) Clear() {
	for i := range m.shards {
		s := &m.shards[i]
		s.clear()
	}
}

// Contains returns true if a key exists.
func (m *shardedMap[K, V]) Contains(key K) bool {
	s := m.shard(key)
	return s.contains(key)
}

// Get returns a value by key, or false.
func (m *shardedMap[K, V]) Get(key K) (V, bool) {
	s := m.shard(key)
	return s.get(key)
}

// GetOrSet returns a value by key and true, or sets a value and false.
func (m *shardedMap[K, V]) GetOrSet(key K, value V) (V, bool) {
	s := m.shard(key)
	return s.getOrSet(key, value)
}

// Delete deletes a key value, and returns the previous value.
func (m *shardedMap[K, V]) Delete(key K) (V, bool) {
	s := m.shard(key)
	return s.delete(key)
}

// LockMap exclusively locks the map.
func (m *shardedMap[K, V]) LockMap() LockedMap[K, V] {
	panic("implement me")
}

// Set sets a value for a key.
func (m *shardedMap[K, V]) Set(key K, value V) {
	s := m.shard(key)
	s.set(key, value)
}

// SetAbsent sets a key value if absent, returns true if set.
func (m *shardedMap[K, V]) SetAbsent(key K, value V) bool {
	s := m.shard(key)
	return s.setAbsent(key, value)
}

// Swap swaps a key value and returns the previous value.
func (m *shardedMap[K, V]) Swap(key K, value V) (V, bool) {
	s := m.shard(key)
	return s.swap(key, value)
}

// Range iterates over all key-value pairs, locks shards during iteration.
func (m *shardedMap[K, V]) Range(fn func(K, V) bool) {
	for i := range m.shards {
		s := &m.shards[i]
		ok := s.range_(fn)
		if !ok {
			return
		}
	}
}

// private

func (m *shardedMap[K, V]) shard(key K) *shardedMapShard[K, V] {
	h := hashing.Hash(key)
	i := int(h) % len(m.shards)
	return &m.shards[i]
}
