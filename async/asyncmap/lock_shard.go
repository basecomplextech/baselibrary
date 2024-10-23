// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asyncmap

import (
	"sync"

	"github.com/basecomplextech/baselibrary/opt"
	"github.com/basecomplextech/baselibrary/pools"
)

type lockShard[K comparable] struct {
	pool pools.Pool[*lockItem[K]]

	mu   sync.Mutex
	item opt.Opt[*lockItem[K]]
	more opt.Opt[map[K]*lockItem[K]]
}

func newLockShard[K comparable](pool pools.Pool[*lockItem[K]]) lockShard[K] {
	return lockShard[K]{pool: pool}
}

// get returns an item by a key, increments its refs, or inserts a new one.
func (s *lockShard[K]) get(key K) *lockItem[K] {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Maybe return existing
	if m, ok := s.item.Unwrap(); ok {
		if m.key == key {
			m.refs++
			return m
		}
	}
	if more, ok := s.more.Unwrap(); ok {
		m, ok := more[key]
		if ok {
			m.refs++
			return m
		}
	}

	// Make new item
	m := s.pool.New()
	m.shard = s
	m.key = key
	m.refs = 1

	// Maybe set single item
	if !s.item.Valid {
		s.item.Set(m)
		return m
	}

	// Othewise add to more items
	more, ok := s.more.Unwrap()
	if !ok {
		more = make(map[K]*lockItem[K])
		s.more.Set(more)
	}
	more[key] = m
	return m
}

// getNoRetain returns an item without incrementing refs, used in tests.
func (s *lockShard[K]) getNoRetain(key K) (*lockItem[K], bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Maybe return existing
	if m, ok := s.item.Unwrap(); ok {
		if m.key == key {
			return m, true
		}
	}
	if more, ok := s.more.Unwrap(); ok {
		m, ok := more[key]
		if ok {
			return m, true
		}
	}
	return nil, false
}

// free releases and removes an item when refs reach zero.
func (s *lockShard[K]) free(m *lockItem[K]) (deleted bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Decrement refs
	if m.refs <= 0 {
		panic("free of freed key lock")
	}
	m.refs--
	if m.refs > 0 {
		return false
	}

	// Delete item
	if m1, ok := s.item.Unwrap(); ok {
		if m1 == m {
			s.item.Unset()
			return true
		}
	}
	if more, ok := s.more.Unwrap(); ok {
		delete(more, m.key)
		return true
	}
	return true
}

// contains returns true if a key exists.
func (s *lockShard[K]) contains(key K) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if m, ok := s.item.Unwrap(); ok {
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

// containsLocked returns true if a key exists, must be called with lock held.
func (s *lockShard[K]) containsLocked(key K) bool {
	if m, ok := s.item.Unwrap(); ok {
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

// rangeLocked calls a function for each key, must be called with lock held.
func (s *lockShard[K]) rangeLocked(f func(key K) bool) {
	if m, ok := s.item.Unwrap(); ok {
		if !f(m.key) {
			return
		}
	}
	if more, ok := s.more.Unwrap(); ok {
		for key := range more {
			if !f(key) {
				break
			}
		}
	}
}
