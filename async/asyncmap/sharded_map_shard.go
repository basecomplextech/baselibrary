// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asyncmap

import (
	"sync"

	"github.com/basecomplextech/baselibrary/opt"
)

type shardedMapShard[K comparable, V any] struct {
	mu sync.RWMutex

	entry opt.Opt[shardedMapEntry[K, V]]
	more  opt.Opt[map[K]V]

	_ [192]byte // cache line padding
}

type shardedMapEntry[K comparable, V any] struct {
	key   K
	value V
}

func (s *shardedMapShard[K, V]) len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	n := 0
	if s.entry.Valid {
		n++
	}
	if more, ok := s.more.Unwrap(); ok {
		n += len(more)
	}
	return n
}

func (s *shardedMapShard[K, V]) clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.entry.Clear()

	if more, ok := s.more.Unwrap(); ok {
		clear(more)
	}
}

func (s *shardedMapShard[K, V]) contains(key K) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s._contains(key)
}

func (s *shardedMapShard[K, V]) get(key K) (v V, _ bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if m, ok := s.entry.Unwrap(); ok {
		if m.key == key {
			return m.value, true
		}
	}
	if more, ok := s.more.Unwrap(); ok {
		v, ok := more[key]
		return v, ok
	}
	return v, false
}

func (s *shardedMapShard[K, V]) getOrSet(key K, value V) (v V, ok bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if m, ok := s.entry.Unwrap(); ok {
		if m.key == key {
			return m.value, true
		}
	}
	if more, ok := s.more.Unwrap(); ok {
		if v, ok := more[key]; ok {
			return v, true
		}
	}

	if !s.entry.Valid {
		e := shardedMapEntry[K, V]{key: key, value: value}
		s.entry.Set(e)
		return v, false
	}

	more, ok := s.more.Unwrap()
	if !ok {
		more = make(map[K]V)
		s.more.Set(more)
	}
	more[key] = value
	return value, false
}

func (s *shardedMapShard[K, V]) delete(key K) (v V, _ bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if m, ok := s.entry.Unwrap(); ok {
		if m.key == key {
			s.entry.Clear()
			return m.value, true
		}
	}
	if more, ok := s.more.Unwrap(); ok {
		v, ok = more[key]
		if ok {
			delete(more, key)
		}
		return v, ok
	}
	return v, false
}

func (s *shardedMapShard[K, V]) set(key K, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s._set(key, value)
}

func (s *shardedMapShard[K, V]) setAbsent(key K, value V) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	ok := s._contains(key)
	if ok {
		return false
	}

	s._set(key, value)
	return true
}

func (s *shardedMapShard[K, V]) swap(key K, value V) (v V, _ bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if m, ok := s.entry.Unwrap(); ok {
		if m.key == key {
			v1 := m.value
			m.value = value
			return v1, true
		}
	}
	if more, ok := s.more.Unwrap(); ok {
		if v1, ok := more[key]; ok {
			more[key] = value
			return v1, true
		}
	}
	return v, false
}

func (s *shardedMapShard[K, V]) range_(fn func(K, V) bool) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if m, ok := s.entry.Unwrap(); ok {
		if !fn(m.key, m.value) {
			return false
		}
	}
	if more, ok := s.more.Unwrap(); ok {
		for key, value := range more {
			if !fn(key, value) {
				return false
			}
		}
	}
	return true
}

// private

func (s *shardedMapShard[K, V]) _contains(key K) bool {
	if m, ok := s.entry.Unwrap(); ok {
		if m.key == key {
			return true
		}
	}
	if more, ok := s.more.Unwrap(); ok {
		_, ok := more[key]
		return ok
	}
	return false
}

func (s *shardedMapShard[K, V]) _set(key K, value V) {
	if !s.entry.Valid {
		e := shardedMapEntry[K, V]{key: key, value: value}
		s.entry.Set(e)
		return
	}

	more, ok := s.more.Unwrap()
	if !ok {
		more = make(map[K]V)
		s.more.Set(more)
	}
	more[key] = value
}
