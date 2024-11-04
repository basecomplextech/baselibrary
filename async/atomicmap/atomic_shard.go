// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package atomics

import (
	"sync"
	"sync/atomic"

	"github.com/basecomplextech/baselibrary/pools"
)

type atomicShard[K comparable, V any] struct {
	pool pools.Pool[*atomicEntry[K, V]]

	wmu   sync.RWMutex
	state atomic.Pointer[atomicState[K, V]]
}

func (s *atomicShard[K, V]) init(size int, pool pools.Pool[*atomicEntry[K, V]]) {
	state := newAtomicState(size, pool)
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

func (s *atomicShard[K, V]) contains(key K) bool {
	state := s.state.Load()
	return state.contains(key)
}

func (s *atomicShard[K, V]) get(key K) (V, bool) {
	state := s.state.Load()
	return state.get(key)
}

func (s *atomicShard[K, V]) getOrSet(key K, value V) (V, bool) {
	resize := false
	v, ok := s._getOrSet(key, value, &resize)

	if resize {
		s._resize()
	}
	return v, ok
}

func (s *atomicShard[K, V]) delete(key K) {
	s.wmu.RLock()
	defer s.wmu.RUnlock()

	state := s.state.Load()
	state.delete(key)
}

func (s *atomicShard[K, V]) pop(key K) (V, bool) {
	s.wmu.RLock()
	defer s.wmu.RUnlock()

	state := s.state.Load()
	return state.pop(key)
}

func (s *atomicShard[K, V]) set(key K, value V) {
	resize := false
	s._set(key, value, &resize)

	if resize {
		s._resize()
	}
}

func (s *atomicShard[K, V]) swap(key K, value V) (V, bool) {
	resize := false
	v, ok := s._swap(key, value, &resize)

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

func (s *atomicShard[K, V]) _getOrSet(key K, value V, resize *bool) (V, bool) {
	s.wmu.RLock()
	defer s.wmu.RUnlock()

	state := s.state.Load()
	v, ok := state.getOrSet(key, value)

	*resize = state.count >= int64(state.threshold)
	return v, ok
}

func (s *atomicShard[K, V]) _set(key K, value V, resize *bool) {
	s.wmu.RLock()
	defer s.wmu.RUnlock()

	state := s.state.Load()
	state.set(key, value)

	*resize = state.count >= int64(state.threshold)
}

func (s *atomicShard[K, V]) _swap(key K, value V, resize *bool) (V, bool) {
	s.wmu.RLock()
	defer s.wmu.RUnlock()

	state := s.state.Load()
	v, ok := state.swap(key, value)

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

	next := resizeAtomicState(state)
	s.state.Store(next)
}
