// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

import (
	"sync"

	"github.com/basecomplextech/baselibrary/opt"
)

type shardedVarShard[T any] struct {
	mu  sync.RWMutex
	ref R[T]
	_   [216]byte
}

func (s *shardedVarShard[T]) acquire() (R[T], bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.ref == nil {
		return nil, false
	}

	ref := s.ref
	ref.Retain()
	return ref, true
}

func (s *shardedVarShard[T]) set(ref R[T]) {
	s.mu.Lock()
	prev := s.ref
	s.ref = ref
	s.mu.Unlock()

	if prev != nil {
		prev.Release()
	}
}

func (s *shardedVarShard[T]) unset() {
	s.mu.Lock()
	prev := s.ref
	s.ref = nil
	s.mu.Unlock()

	if prev != nil {
		prev.Release()
	}
}

func (s *shardedVarShard[T]) unwrap() opt.Opt[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.ref == nil {
		return opt.None[T]()
	}

	value := s.ref.Unwrap()
	return opt.New(value)
}

func (s *shardedVarShard[T]) unwrapRef() opt.Opt[R[T]] {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.ref == nil {
		return opt.None[R[T]]()
	}

	return opt.New[R[T]](s.ref)
}
