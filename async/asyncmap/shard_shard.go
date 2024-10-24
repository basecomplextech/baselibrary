// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asyncmap

import (
	"sync"

	"github.com/basecomplextech/baselibrary/opt"
)

type mapShard[K comparable, V any] struct {
	mu    sync.Mutex
	entry opt.Opt[mapEntry[K, V]]
	more  opt.Opt[map[K]V]
}

type mapEntry[K comparable, V any] struct {
	key   K
	value V
}

func newMapShard[K comparable, V any]() mapShard[K, V] {
	return mapShard[K, V]{}
}

func (s *mapShard[K, V]) len() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	n := 0
	if s.entry.Valid {
		n++
	}
	if more, ok := s.more.Unwrap(); ok {
		n += len(more)
	}
	return n
}

func (s *mapShard[K, V]) clear() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.entry.Unset()

	if more, ok := s.more.Unwrap(); ok {
		clear(more)
	}
}

func (s *mapShard[K, V]) contains(key K) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

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

func (s *mapShard[K, V]) get(key K) (v V, _ bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

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

func (s *mapShard[K, V]) getOrSet(key K, value V) (v V, set bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if m, ok := s.entry.Unwrap(); ok {
		if m.key == key {
			return m.value, false
		}
	}
	if more, ok := s.more.Unwrap(); ok {
		if v, ok := more[key]; ok {
			return v, false
		}
	}

	if !s.entry.Valid {
		e := mapEntry[K, V]{key: key, value: value}
		s.entry.Set(e)
		return v, true
	}

	more, ok := s.more.Unwrap()
	if !ok {
		more = make(map[K]V)
		s.more.Set(more)
	}
	more[key] = value
	return value, true
}

func (s *mapShard[K, V]) delete(key K) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if m, ok := s.entry.Unwrap(); ok {
		if m.key == key {
			s.entry.Unset()
			return
		}
	}
	if more, ok := s.more.Unwrap(); ok {
		delete(more, key)
	}
}

func (s *mapShard[K, V]) pop(key K) (v V, _ bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if m, ok := s.entry.Unwrap(); ok {
		if m.key == key {
			s.entry.Unset()
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

func (s *mapShard[K, V]) set(key K, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.entry.Valid {
		e := mapEntry[K, V]{key: key, value: value}
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

func (s *mapShard[K, V]) swap(key K, value V) (v V, _ bool) {
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

func (s *mapShard[K, V]) range_(fn func(K, V) bool) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

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
