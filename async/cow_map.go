// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import (
	"maps"
	"runtime"
	"sync"
	"sync/atomic"

	"github.com/basecomplextech/baselibrary/internal/hashing"
)

// CopyOnWriteMap is a concurrent map which copies the whole map on any write.
//
// The map is optimized for completely non-blocking reads and very infrequent writes.
type CopyOnWriteMap[K comparable, V any] interface {
	Map[K, V]
}

// NewCopyOnWriteMap returns a new copy-on-write map.
func NewCopyOnWriteMap[K comparable, V any]() Map[K, V] {
	return newCowMap[K, V]()
}

// internal

var _ CopyOnWriteMap[int, int] = (*cowMap[int, int])(nil)

type cowMap[K comparable, V any] struct {
	shards []cowMapShard[K, V]
}

func newCowMap[K comparable, V any]() *cowMap[K, V] {
	num := runtime.NumCPU() * 2

	shards := make([]cowMapShard[K, V], num)
	for i := range shards {
		shards[i].init()
	}

	return &cowMap[K, V]{shards: shards}
}

// Len returns the number of keys.
func (m *cowMap[K, V]) Len() int {
	var n int
	for i := range m.shards {
		s := &m.shards[i]
		n += s.len()
	}
	return n
}

// Clear deletes all items.
func (m *cowMap[K, V]) Clear() {
	for i := range m.shards {
		s := &m.shards[i]
		s.clear()
	}
}

// Contains returns true if a key exists.
func (m *cowMap[K, V]) Contains(key K) bool {
	s := m.shard(key)
	return s.contains(key)
}

// Get returns a value by key, or false.
func (m *cowMap[K, V]) Get(key K) (V, bool) {
	s := m.shard(key)
	return s.get(key)
}

// GetOrSet returns a value by key, or sets a value if it does not exist.
func (m *cowMap[K, V]) GetOrSet(key K, value V) (_ V, set bool) {
	s := m.shard(key)
	return s.getOrSet(key, value)
}

// Pop deletes and returns a value by key, or false.
func (m *cowMap[K, V]) Pop(key K) (V, bool) {
	s := m.shard(key)
	return s.pop(key)
}

// Set sets a value for a key.
func (m *cowMap[K, V]) Set(key K, value V) {
	s := m.shard(key)
	s.set(key, value)
}

// Delete deletes a value by key.
func (m *cowMap[K, V]) Delete(key K) {
	s := m.shard(key)
	s.delete(key)
}

// Range iterates over all key-value pairs, locks shards during iteration.
func (m *cowMap[K, V]) Range(fn func(K, V)) {
	for i := range m.shards {
		s := &m.shards[i]
		s.range_(fn)
	}
}

// private

func (m *cowMap[K, V]) shard(key K) *cowMapShard[K, V] {
	index := hashing.Shard(key, len(m.shards))
	return &m.shards[index]
}

// shard

type cowMapShard[K comparable, V any] struct {
	current atomic.Pointer[cowMapVersion[K, V]]
	wmu     sync.Mutex

	_ [240]byte // cache line padding
}

func (s *cowMapShard[K, V]) init() {
	v := newCowVersion[K, V]()
	s.current.Store(v)
}

func (s *cowMapShard[K, V]) len() int {
	v := s.current.Load()
	return len(v.items)
}

func (s *cowMapShard[K, V]) clear() {
	s.wmu.Lock()
	defer s.wmu.Unlock()

	v1 := newCowVersion[K, V]()
	s.current.Store(v1)
}

func (s *cowMapShard[K, V]) contains(key K) bool {
	v := s.current.Load()
	_, ok := v.items[key]
	return ok
}

func (s *cowMapShard[K, V]) get(key K) (V, bool) {
	v := s.current.Load()
	value, ok := v.items[key]
	return value, ok
}

func (s *cowMapShard[K, V]) getOrSet(key K, value V) (V, bool) {
	value1, ok := s.get(key)
	if ok {
		return value1, false
	}

	// Slow path
	s.wmu.Lock()
	defer s.wmu.Unlock()

	// Check again
	v := s.current.Load()
	value1, ok = v.items[key]
	if ok {
		return value1, false
	}

	// Set new value
	v1 := v.clone()
	v1.items[key] = value
	s.current.Store(v1)
	return value, true
}

func (s *cowMapShard[K, V]) pop(key K) (V, bool) {
	s.wmu.Lock()
	defer s.wmu.Unlock()

	v := s.current.Load()
	value, ok := v.items[key]
	if !ok {
		return value, false
	}

	v1 := v.clone()
	delete(v1.items, key)
	s.current.Store(v1)
	return value, true
}

func (s *cowMapShard[K, V]) set(key K, value V) {
	s.wmu.Lock()
	defer s.wmu.Unlock()

	v := s.current.Load()
	v1 := v.clone()
	v1.items[key] = value
	s.current.Store(v1)
}

func (s *cowMapShard[K, V]) delete(key K) {
	s.wmu.Lock()
	defer s.wmu.Unlock()

	v := s.current.Load()
	if _, ok := v.items[key]; !ok {
		return
	}

	v1 := v.clone()
	delete(v1.items, key)
	s.current.Store(v1)
}

func (s *cowMapShard[K, V]) range_(fn func(K, V)) {
	v := s.current.Load()

	for key, value := range v.items {
		fn(key, value)
	}
}

// version

type cowMapVersion[K comparable, V any] struct {
	items map[K]V
}

func newCowVersion[K comparable, V any]() *cowMapVersion[K, V] {
	return &cowMapVersion[K, V]{
		items: make(map[K]V),
	}
}

func (v *cowMapVersion[K, V]) clone() *cowMapVersion[K, V] {
	return &cowMapVersion[K, V]{
		items: maps.Clone(v.items),
	}
}
