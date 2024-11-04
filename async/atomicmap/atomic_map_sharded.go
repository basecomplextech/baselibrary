// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package atomics

import (
	"math/bits"
	"runtime"
	"unsafe"

	"github.com/basecomplextech/baselibrary/async/asyncmap"
	"github.com/basecomplextech/baselibrary/internal/hashing"
	"github.com/basecomplextech/baselibrary/pools"
)

var _ asyncmap.Map[int, int] = (*atomicShardedMap[int, int])(nil)

type atomicShardedMap[K comparable, V any] struct {
	pool   pools.Pool[*atomicEntry[K, V]]
	shards []atomicShard[K, V] // always power of two

	bitWidth int // shard hash bit width
}

func newAtomicShardedMap[K comparable, V any](size int) *atomicShardedMap[K, V] {
	pool := newAtomicPool[K, V]()

	// CPUs and cache line size
	cpus := runtime.NumCPU()
	cacheLineSize := 256

	// Calculate number of shards, round to power of two
	shardSize := unsafe.Sizeof(atomicShard[K, V]{})
	shardNum := (cpus * cacheLineSize) / int(shardSize)
	shardNum = roundToPowerOfTwo(shardNum)

	// Make shards
	shards := make([]atomicShard[K, V], shardNum)
	for i := range shards {
		size1 := size / shardNum
		shards[i].init(size1, pool)
	}

	// Calculate shard hash bit width
	bitWidth := bitWidth(shardNum - 1)

	// Return map
	return &atomicShardedMap[K, V]{
		pool:     pool,
		shards:   shards,
		bitWidth: bitWidth,
	}
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
	s, _ := m.shard(key)
	return s.contains(key)
}

// Get returns a value by key, or false.
func (m *atomicShardedMap[K, V]) Get(key K) (v V, _ bool) {
	s, _ := m.shard(key)
	return s.get(key)
}

// GetOrSet returns a value by key, or sets a value if it does not exist.
func (m *atomicShardedMap[K, V]) GetOrSet(key K, value V) (_ V, set bool) {
	s, _ := m.shard(key)
	return s.getOrSet(key, value)
}

// Delete deletes a value by key.
func (m *atomicShardedMap[K, V]) Delete(key K) {
	s, _ := m.shard(key)
	s.delete(key)
}

// Pop deletes and returns a value by key, or false.
func (m *atomicShardedMap[K, V]) Pop(key K) (v V, _ bool) {
	s, _ := m.shard(key)
	return s.pop(key)
}

// Set sets a value for a key.
func (m *atomicShardedMap[K, V]) Set(key K, value V) {
	s, _ := m.shard(key)
	s.set(key, value)
}

// Swap swaps a key value and returns the previous value.
func (m *atomicShardedMap[K, V]) Swap(key K, value V) (v V, _ bool) {
	s, _ := m.shard(key)
	return s.swap(key, value)
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

func (m *atomicShardedMap[K, V]) shard(key K) (_ *atomicShard[K, V], h1 uint32) {
	h := hashing.Hash(key)
	h1 = m.shardHash(h)

	i := h % uint32(len(m.shards))
	return &m.shards[i], h1
}

func (m *atomicShardedMap[K, V]) shardHash(h uint32) uint32 {
	if int(h) < len(m.shards) {
		return h
	}
	return h >> m.bitWidth
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
