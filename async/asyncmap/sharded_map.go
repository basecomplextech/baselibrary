// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asyncmap

import (
	"runtime"

	"github.com/basecomplextech/baselibrary/internal/hashing"
)

// ShardedMap is a map which uses multiple hash maps each guarded by a separate mutex.
//
// This map is optimized for read-write operations.
// Use [SyncMap] if you need a map optimized mostly for read operations.
// Use [AtomicMap] or [AtomicShardedMap] if you need a map optimized for read/write
// operations with non-blocking reads.
//
// # Benchmarks
//
//	cpu: Apple M1 Pro
//	BenchmarkShardedMap_Read-10                            	30489063	        37.41 ns/op	        26.73 mops	       0 B/op	       0 allocs/op
//	BenchmarkShardedMap_Read_Parallel-10                   	81285898	        14.79 ns/op	        67.60 mops	       0 B/op	       0 allocs/op
//	BenchmarkShardedMap_Write-10                           	28579593	        42.33 ns/op	        23.62 mops	       0 B/op	       0 allocs/op
//	BenchmarkShardedMap_Write_Parallel-10                  	21006073	        55.74 ns/op	        17.94 mops	       0 B/op	       0 allocs/op
//	BenchmarkShardedMap_Read_Write_Parallel-10             	19586400	        60.62 ns/op	         1.07 rmops	        16.50 wmops	       0 B/op	       0 allocs/op
//	BenchmarkShardedMap_Read_Parallel_Write_Parallel-10    	 3161337	       416.00 ns/op	        33.10 rmops	         2.40 wmops	       0 B/op	       0 allocs/op
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
	cpus := runtime.NumCPU() * 4

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

// Delete deletes a value by key.
func (m *shardedMap[K, V]) Delete(key K) {
	s := m.shard(key)
	s.delete(key)
}

// Pop deletes and returns a value by key, or false.
func (m *shardedMap[K, V]) Pop(key K) (V, bool) {
	s := m.shard(key)
	return s.pop(key)
}

// Set sets a value for a key.
func (m *shardedMap[K, V]) Set(key K, value V) {
	s := m.shard(key)
	s.set(key, value)
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
