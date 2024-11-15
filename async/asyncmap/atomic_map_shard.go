// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asyncmap

import (
	"sync"
	"sync/atomic"
)

type atomicMapShard[K comparable, V any] struct {
	m *atomicShardedMap[K, V]

	wmu   sync.RWMutex
	state atomic.Pointer[atomicMapState[K, V]]

	_ [216]byte // cache line padding
}

func (s *atomicMapShard[K, V]) init(m *atomicShardedMap[K, V], size int) {
	state := newAtomicMapState(size, m.pool)

	s.m = m
	s.state.Store(state)
}

func (s *atomicMapShard[K, V]) len() int {
	state := s.state.Load()
	return state.len()
}

func (s *atomicMapShard[K, V]) clear() {
	s.wmu.Lock()
	defer s.wmu.Unlock()

	state := s.state.Load()
	next := emptyAtomicMapState(state)

	s.state.Store(next)
}

func (s *atomicMapShard[K, V]) contains(h uint32, key K) bool {
	state := s.state.Load()
	return state.contains(h, key)
}

func (s *atomicMapShard[K, V]) get(h uint32, key K) (V, bool) {
	state := s.state.Load()
	return state.get(h, key)
}

func (s *atomicMapShard[K, V]) getOrSet(h uint32, key K, value V) (V, bool) {
	resize := false
	v, ok := s._getOrSet(h, key, value, &resize)

	if resize {
		s._resize()
	}
	return v, ok
}

func (s *atomicMapShard[K, V]) delete(h uint32, key K) (V, bool) {
	s.wmu.RLock()
	defer s.wmu.RUnlock()

	state := s.state.Load()
	return state.delete(h, key)
}

func (s *atomicMapShard[K, V]) set(h uint32, key K, value V) {
	resize := false
	s._set(h, key, value, &resize)

	if resize {
		s._resize()
	}
}

func (s *atomicMapShard[K, V]) setAbsent(h uint32, key K, value V) bool {
	resize := false
	ok := s._setAbsent(h, key, value, &resize)

	if resize {
		s._resize()
	}
	return ok
}

func (s *atomicMapShard[K, V]) swap(h uint32, key K, value V) (V, bool) {
	resize := false
	v, ok := s._swap(h, key, value, &resize)

	if resize {
		s._resize()
	}
	return v, ok
}

func (s *atomicMapShard[K, V]) range_(fn func(K, V) bool) bool {
	state := s.state.Load()
	return state.range_(fn)
}

// private

func (s *atomicMapShard[K, V]) _getOrSet(h uint32, key K, value V, resize *bool) (V, bool) {
	s.wmu.RLock()
	defer s.wmu.RUnlock()

	state := s.state.Load()
	v, ok := state.getOrSet(h, key, value)

	n := state.len()
	*resize = n >= state.threshold
	return v, ok
}

func (s *atomicMapShard[K, V]) _set(h uint32, key K, value V, resize *bool) {
	s.wmu.RLock()
	defer s.wmu.RUnlock()

	state := s.state.Load()
	state.set(h, key, value)

	n := state.len()
	*resize = n >= state.threshold
}

func (s *atomicMapShard[K, V]) _setAbsent(h uint32, key K, value V, resize *bool) bool {
	s.wmu.RLock()
	defer s.wmu.RUnlock()

	state := s.state.Load()
	ok := state.setAbsent(h, key, value)

	n := state.len()
	*resize = n >= state.threshold
	return ok
}

func (s *atomicMapShard[K, V]) _swap(h uint32, key K, value V, resize *bool) (V, bool) {
	s.wmu.RLock()
	defer s.wmu.RUnlock()

	state := s.state.Load()
	v, ok := state.swap(h, key, value)

	n := state.len()
	*resize = n >= state.threshold
	return v, ok
}

// resize

func (s *atomicMapShard[K, V]) _resize() {
	s.wmu.Lock()
	defer s.wmu.Unlock()

	state := s.state.Load()
	n := state.len()
	if n < state.threshold {
		return
	}

	// Double buckets
	size := len(state.buckets) * 2
	next := newAtomicMapState(size, s.m.pool)

	// Copy all items
	state.rangeLocked(func(k K, v V) bool {
		_, h1 := s.m.hashes(k)
		next.set(h1, k, v)
		return true
	})

	// Replace state
	s.state.Store(next)
}
