// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package atomics

import (
	"sync"
	"sync/atomic"
)

type atomicShard[K comparable, V any] struct {
	m *atomicShardedMap[K, V]

	wmu   sync.RWMutex
	state atomic.Pointer[atomicState[K, V]]

	_ [216]byte
}

func (s *atomicShard[K, V]) init(m *atomicShardedMap[K, V], size int) {
	state := newAtomicState(size, m.pool)

	s.m = m
	s.state.Store(state)
}

func (s *atomicShard[K, V]) len() int {
	state := s.state.Load()
	return int(state.count)
}

func (s *atomicShard[K, V]) clear() {
	s.wmu.Lock()
	defer s.wmu.Unlock()

	state := s.state.Load()
	next := clearAtomicState(state)

	s.state.Store(next)
}

func (s *atomicShard[K, V]) contains(h uint32, key K) bool {
	state := s.state.Load()
	return state.contains(h, key)
}

func (s *atomicShard[K, V]) get(h uint32, key K) (V, bool) {
	state := s.state.Load()
	return state.get(h, key)
}

func (s *atomicShard[K, V]) getOrSet(h uint32, key K, value V) (V, bool) {
	resize := false
	v, ok := s._getOrSet(h, key, value, &resize)

	if resize {
		s._resize()
	}
	return v, ok
}

func (s *atomicShard[K, V]) delete(h uint32, key K) {
	s.wmu.RLock()
	defer s.wmu.RUnlock()

	state := s.state.Load()
	state.delete(h, key)
}

func (s *atomicShard[K, V]) pop(h uint32, key K) (V, bool) {
	s.wmu.RLock()
	defer s.wmu.RUnlock()

	state := s.state.Load()
	return state.pop(h, key)
}

func (s *atomicShard[K, V]) set(h uint32, key K, value V) {
	resize := false
	s._set(h, key, value, &resize)

	if resize {
		s._resize()
	}
}

func (s *atomicShard[K, V]) swap(h uint32, key K, value V) (V, bool) {
	resize := false
	v, ok := s._swap(h, key, value, &resize)

	if resize {
		s._resize()
	}
	return v, ok
}

func (s *atomicShard[K, V]) range_(fn func(K, V) bool) bool {
	state := s.state.Load()
	return state.range_(fn)
}

// private

func (s *atomicShard[K, V]) _getOrSet(h uint32, key K, value V, resize *bool) (V, bool) {
	s.wmu.RLock()
	defer s.wmu.RUnlock()

	state := s.state.Load()
	v, ok := state.getOrSet(h, key, value)

	*resize = state.count >= int64(state.threshold)
	return v, ok
}

func (s *atomicShard[K, V]) _set(h uint32, key K, value V, resize *bool) {
	s.wmu.RLock()
	defer s.wmu.RUnlock()

	state := s.state.Load()
	state.set(h, key, value)

	*resize = state.count >= int64(state.threshold)
}

func (s *atomicShard[K, V]) _swap(h uint32, key K, value V, resize *bool) (V, bool) {
	s.wmu.RLock()
	defer s.wmu.RUnlock()

	state := s.state.Load()
	v, ok := state.swap(h, key, value)

	*resize = state.count >= int64(state.threshold)
	return v, ok
}

// resize

func (s *atomicShard[K, V]) _resize() {
	s.wmu.Lock()
	defer s.wmu.Unlock()

	state := s.state.Load()
	if state.count < int64(state.threshold) {
		return
	}

	// Double buckets
	size := len(state.buckets) * 2
	next := newAtomicState(size, s.m.pool)

	// Copy all items
	state.rangeLocked(func(k K, v V) bool {
		_, h1 := s.m.hashes(k)
		next.set(h1, k, v)
		return true
	})

	// Replace state
	s.state.Store(next)
}
