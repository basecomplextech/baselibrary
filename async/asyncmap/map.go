// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asyncmap

import (
	"runtime"
	"sync"
	"unsafe"

	"github.com/basecomplextech/baselibrary/internal/hashing"
)

// Map is a sharded hash map, which uses a lock per shard.
// The number of shards is equal to the number of CPU cores.
//
// This map is optimized for read-write operations.
// Maybe use ConcurrentMap if you need a map optimized mostly for read operations.
type Map[K comparable, V any] interface {
	// Len returns the number of keys.
	Len() int

	// Clear deletes all items.
	Clear()

	// Contains returns true if a key exists.
	Contains(key K) bool

	// Get returns a value by key, or false.
	Get(key K) (V, bool)

	// GetOrSet returns a value by key, or sets a value if it does not exist.
	GetOrSet(key K, value V) (_ V, set bool)

	// Pop deletes and returns a value by key, or false.
	Pop(key K) (V, bool)

	// Set sets a value for a key.
	Set(key K, value V)

	// Delete deletes a value by key.
	Delete(key K)

	// Range iterates over all key-value pairs, locks shards during iteration.
	Range(fn func(K, V))
}

// NewMap returns a new sharded map.
func NewMap[K comparable, V any]() Map[K, V] {
	return newShardedMap[K, V]()
}

// internal

var _ Map[int, int] = &shardedMap[int, int]{}

type shardedMap[K comparable, V any] struct {
	shards []mapShard[K, V]
}

func newShardedMap[K comparable, V any]() *shardedMap[K, V] {
	cpus := runtime.NumCPU()
	cpuLines := 16
	lineSize := 256

	size := unsafe.Sizeof(mapShard[K, V]{})
	total := cpus * cpuLines * lineSize
	n := int(total / int(size))

	// n := runtime.NumCPU() * 4
	shards := make([]mapShard[K, V], n)
	for i := range shards {
		shards[i] = newMapShard[K, V]()
	}

	return &shardedMap[K, V]{shards: shards}
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

// GetOrSet returns a value by key, or sets a value if it does not exist.
func (m *shardedMap[K, V]) GetOrSet(key K, value V) (_ V, set bool) {
	s := m.shard(key)
	return s.getOrSet(key, value)
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

// Delete deletes a value by key.
func (m *shardedMap[K, V]) Delete(key K) {
	s := m.shard(key)
	s.delete(key)
}

// Range iterates over all key-value pairs, locks shards during iteration.
func (m *shardedMap[K, V]) Range(fn func(K, V)) {
	for i := range m.shards {
		s := &m.shards[i]
		s.range_(fn)
	}
}

// private

func (m *shardedMap[K, V]) shard(key K) *mapShard[K, V] {
	index := hashing.Shard(key, len(m.shards))
	return &m.shards[index]
}

// shard

type mapShard[K comparable, V any] struct {
	mu    sync.RWMutex
	items map[K]V

	// _ [224]byte
}

func newMapShard[K comparable, V any]() mapShard[K, V] {
	return mapShard[K, V]{
		items: make(map[K]V),
	}
}

func (s *mapShard[K, V]) len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.items)
}

func (s *mapShard[K, V]) clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	clear(s.items)
}

func (s *mapShard[K, V]) contains(key K) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.items[key]
	return ok
}

func (s *mapShard[K, V]) get(key K) (V, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, ok := s.items[key]
	return value, ok
}

func (s *mapShard[K, V]) getOrSet(key K, value V) (V, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if v, ok := s.items[key]; ok {
		return v, false
	}

	s.items[key] = value
	return value, true
}

func (s *mapShard[K, V]) pop(key K) (V, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, ok := s.items[key]
	if !ok {
		return value, false
	}

	delete(s.items, key)
	return value, true
}

func (s *mapShard[K, V]) set(key K, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.items[key] = value
}

func (s *mapShard[K, V]) delete(key K) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.items, key)
}

func (s *mapShard[K, V]) range_(fn func(K, V)) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for key, value := range s.items {
		fn(key, value)
	}
}
