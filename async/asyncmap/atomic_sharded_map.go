// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asyncmap

import (
	"math/bits"
	"runtime"

	"github.com/basecomplextech/baselibrary/internal/hashing"
	"github.com/basecomplextech/baselibrary/pools"
)

// AtomicShardedMap is a goroutine-safe hash map based on atomic operations,
// multiple shards to reduce contention, and hierarchical hashing.
//
// Readers are non-blocking, writers use a mutex per bucket, and a resize mutex per shard.
// See [AtomicMap] for the implementation details.
//
// Benchmarks:
//
//	cpu: Apple M1 Pro
//	BenchmarkAtomicShardedMap_Read-10                            	67143218	        17.81 ns/op	        56.15 mops	       0 B/op	       0 allocs/op
//	BenchmarkAtomicShardedMap_Read_Parallel-10                   	66365866	        19.75 ns/op	        50.63 mops	       0 B/op	       0 allocs/op
//	BenchmarkAtomicShardedMap_Write-10                           	26752994	        44.23 ns/op	        22.61 mops	       0 B/op	       0 allocs/op
//	BenchmarkAtomicShardedMap_Write_Parallel-10                  	24188610	        43.98 ns/op	        22.74 mops	       0 B/op	       0 allocs/op
//	BenchmarkAtomicShardedMap_Read_Write_Parallel-10             	24252631	        49.88 ns/op	         5.58 rmops	       20.05 wmops	       0 B/op	       0 allocs/op
//	BenchmarkAtomicShardedMap_Read_Parallel_Write_Parallel-10    	 5750632	       237.40 ns/op	        43.42 rmops	        4.21 wmops	       0 B/op	       0 allocs/op
type AtomicShardedMap[K comparable, V any] interface {
	Map[K, V]
}

// internal

var _ Map[int, int] = (*atomicShardedMap[int, int])(nil)

type atomicShardedMap[K comparable, V any] struct {
	pool   pools.Pool[*atomicMapEntry[K, V]]
	shards []atomicMapShard[K, V] // always power of two

	bitWidth int // shard hash bit width
}

func newAtomicShardedMap[K comparable, V any](size int) *atomicShardedMap[K, V] {
	pool := newAtomicPool[K, V]()

	// Calculate shard number, round to power of two
	cpus := runtime.NumCPU()
	shardNum := roundToPowerOfTwo(cpus)
	bitWidth := bitWidth(shardNum - 1)

	// Make map
	m := &atomicShardedMap[K, V]{
		pool:     pool,
		shards:   make([]atomicMapShard[K, V], shardNum),
		bitWidth: bitWidth,
	}

	// Init shards
	for i := range m.shards {
		size1 := size / shardNum
		m.shards[i].init(m, size1)
	}
	return m
}

// Len returns the number of keys.
func (m *atomicShardedMap[K, V]) Len() int {
	n := 0
	for i := range m.shards {
		n += m.shards[i].len()
	}
	return n
}

// Clear deletes all items.
func (m *atomicShardedMap[K, V]) Clear() {
	for i := range m.shards {
		m.shards[i].clear()
	}
}

// Contains returns true if a key exists.
func (m *atomicShardedMap[K, V]) Contains(key K) bool {
	s, h := m.shard(key)
	return s.contains(h, key)
}

// Get returns a value by key, or false.
func (m *atomicShardedMap[K, V]) Get(key K) (v V, _ bool) {
	s, h := m.shard(key)
	return s.get(h, key)
}

// GetOrSet returns a value by key, or sets a value if it does not exist.
func (m *atomicShardedMap[K, V]) GetOrSet(key K, value V) (_ V, set bool) {
	s, h := m.shard(key)
	return s.getOrSet(h, key, value)
}

// Delete deletes a value by key.
func (m *atomicShardedMap[K, V]) Delete(key K) {
	s, h := m.shard(key)
	s.delete(h, key)
}

// Pop deletes and returns a value by key, or false.
func (m *atomicShardedMap[K, V]) Pop(key K) (v V, _ bool) {
	s, h := m.shard(key)
	return s.pop(h, key)
}

// Set sets a value for a key.
func (m *atomicShardedMap[K, V]) Set(key K, value V) {
	s, h := m.shard(key)
	s.set(h, key, value)
}

// Swap swaps a key value and returns the previous value.
func (m *atomicShardedMap[K, V]) Swap(key K, value V) (v V, _ bool) {
	s, h := m.shard(key)
	return s.swap(h, key, value)
}

// Range iterates over all key-value pairs.
// The iteration stops if the function returns false.
func (m *atomicShardedMap[K, V]) Range(fn func(K, V) bool) {
	for i := range m.shards {
		ok := m.shards[i].range_(fn)
		if !ok {
			return
		}
	}
}

// private

// hashes returns two hierarchical key hashes, one for map and one for shard.
func (m *atomicShardedMap[K, V]) hashes(key K) (uint32, uint32) {
	h := hashing.Hash(key)
	h1 := h

	// Right shift to get shard hash
	//
	// For example, 16 is 10000, 4 bits required to store shard index [0-15].
	// If hash is >= 16, then we can use higher bits for shard hash.
	if int(h1) >= len(m.shards) {
		h1 = h1 >> m.bitWidth
	}
	return h, h1
}

func (m *atomicShardedMap[K, V]) shard(key K) (_ *atomicMapShard[K, V], h1 uint32) {
	h, h1 := m.hashes(key)
	i := h % uint32(len(m.shards))
	return &m.shards[i], h1
}

// util

// bitWidth returns the number of bits required to represent a number.
func bitWidth(n int) int {
	zeros := bits.LeadingZeros64(uint64(n))
	return 64 - zeros
}

// roundToPowerOfTwo rounds a number to the nearest power of two.
//
// See https://stackoverflow.com/questions/466204/rounding-up-to-next-power-of-2
func roundToPowerOfTwo(n int) int {
	n--
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	n++
	return n
}
