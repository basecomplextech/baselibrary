// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

import (
	"sync/atomic"

	"github.com/basecomplextech/baselibrary/opt"
)

type concVarShard[T any] struct {
	ref atomic.Pointer[concValueRef[T]]
	_   [256 - 8]byte // padding
}

func (s *concVarShard[T]) acquire() (R[T], bool) {
	ref := s.ref.Load()
	if ref == nil {
		return nil, false
	}

	if ref.Acquire() {
		return ref, true
	}
	return nil, false
}

func (s *concVarShard[T]) set(ref *concValueRef[T]) {
	prev := s.ref.Swap(ref)
	if prev == nil {
		return
	}

	prev.Release()
}

func (s *concVarShard[T]) unset() {
	prev := s.ref.Swap(nil)
	if prev == nil {
		return
	}

	prev.Release()
}

func (s *concVarShard[T]) unwrap() opt.Opt[T] {
	r := s.ref.Load()
	if r == nil {
		return opt.None[T]()
	}

	value := r.Unwrap()
	return opt.New[T](value)
}

func (s *concVarShard[T]) unwrapRef() opt.Opt[R[T]] {
	r := s.ref.Load()
	if r == nil {
		return opt.None[R[T]]()
	}

	return opt.New[R[T]](r)
}
