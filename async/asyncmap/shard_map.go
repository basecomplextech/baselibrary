// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asyncmap

import (
	"runtime"
	"unsafe"

	"github.com/basecomplextech/baselibrary/internal/hashing"
)

// ShardMap is a map which internally uses multiple hash maps each guarded by a separate mutex.
//
// This map is optimized for read-write operations.
// Use AtomicMap if you need a map optimized mostly for write operations.
type ShardMap[K comparable, V any] interface {
	Map[K, V]
}

// NewShardMap returns a new sharded map.
func NewShardMap[K comparable, V any]() Map[K, V] {
	return newShardMap[K, V]()
}

// internal

var _ Map[int, int] = &shardMap[int, int]{}

type shardMap[K comparable, V any] struct {
	shards []mapShard[K, V]
}

func newShardMap[K comparable, V any]() *shardMap[K, V] {
	cpus := runtime.NumCPU()
	lineSize := 256
	linesPerCPU := 16

	size := unsafe.Sizeof(mapShard[K, V]{})
	total := cpus * linesPerCPU * lineSize
	n := int(total / int(size))

	// n := runtime.NumCPU() * 4
	shards := make([]mapShard[K, V], n)
	for i := range shards {
		shards[i] = newMapShard[K, V]()
	}

	return &shardMap[K, V]{shards: shards}
}

// Len returns the number of keys.
func (m *shardMap[K, V]) Len() int {
	var n int

	for i := range m.shards {
		s := &m.shards[i]
		n += s.len()
	}

	return n
}

// Clear deletes all items.
func (m *shardMap[K, V]) Clear() {
	for i := range m.shards {
		s := &m.shards[i]
		s.clear()
	}
}

// Contains returns true if a key exists.
func (m *shardMap[K, V]) Contains(key K) bool {
	s := m.shard(key)
	return s.contains(key)
}

// Get returns a value by key, or false.
func (m *shardMap[K, V]) Get(key K) (V, bool) {
	s := m.shard(key)
	return s.get(key)
}

// GetOrSet returns a value by key, or sets a value if it does not exist.
func (m *shardMap[K, V]) GetOrSet(key K, value V) (_ V, set bool) {
	s := m.shard(key)
	return s.getOrSet(key, value)
}

// Delete deletes a value by key.
func (m *shardMap[K, V]) Delete(key K) {
	s := m.shard(key)
	s.delete(key)
}

// Pop deletes and returns a value by key, or false.
func (m *shardMap[K, V]) Pop(key K) (V, bool) {
	s := m.shard(key)
	return s.pop(key)
}

// Set sets a value for a key.
func (m *shardMap[K, V]) Set(key K, value V) {
	s := m.shard(key)
	s.set(key, value)
}

// Swap swaps a key value and returns the previous value.
func (m *shardMap[K, V]) Swap(key K, value V) (V, bool) {
	s := m.shard(key)
	return s.swap(key, value)
}

// Range iterates over all key-value pairs, locks shards during iteration.
func (m *shardMap[K, V]) Range(fn func(K, V) bool) {
	for i := range m.shards {
		s := &m.shards[i]
		ok := s.range_(fn)
		if !ok {
			return
		}
	}
}

// private

func (m *shardMap[K, V]) shard(key K) *mapShard[K, V] {
	index := hashing.Shard(key, len(m.shards))
	return &m.shards[index]
}
