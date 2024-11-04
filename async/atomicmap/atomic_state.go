// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package atomics

import (
	"sync/atomic"

	"github.com/basecomplextech/baselibrary/internal/hashing"
	"github.com/basecomplextech/baselibrary/pools"
)

type atomicState[K comparable, V any] struct {
	pool pools.Pool[*atomicEntry[K, V]]

	count     int64 // number of items
	threshold int   // next resize threshold

	buckets []atomicBucket[K, V]
}

func newAtomicState[K comparable, V any](
	size int,
	pool pools.Pool[*atomicEntry[K, V]],
) *atomicState[K, V] {

	size = max(size, atomicMapMinSize)
	threshold := int(float64(size) * atomicMapThreshold)

	return &atomicState[K, V]{
		pool:      pool,
		threshold: threshold,

		buckets: make([]atomicBucket[K, V], size),
	}
}

func clearAtomicState[K comparable, V any](s *atomicState[K, V]) *atomicState[K, V] {
	return &atomicState[K, V]{
		pool:      s.pool,
		threshold: s.threshold,

		buckets: make([]atomicBucket[K, V], len(s.buckets)),
	}
}

func resizeAtomicState[K comparable, V any](s *atomicState[K, V]) *atomicState[K, V] {
	size := len(s.buckets) * 2
	next := newAtomicState(size, s.pool)

	s.rangeLocked(func(k K, v V) bool {
		next.set(k, v)
		return true
	})
	return next
}

// methods

func (s *atomicState[K, V]) contains(key K) bool {
	b := s.bucket(key)
	_, ok := b.get(key, s.pool)
	return ok
}

func (s *atomicState[K, V]) get(key K) (V, bool) {
	b := s.bucket(key)
	return b.get(key, s.pool)
}

func (s *atomicState[K, V]) getOrSet(key K, value V) (V, bool) {
	b := s.bucket(key)
	v, ok := b.getOrSet(key, value, s.pool)
	if !ok {
		atomic.AddInt64(&s.count, 1)
	}
	return v, ok
}

func (s *atomicState[K, V]) delete(key K) {
	b := s.bucket(key)
	_, ok := b.delete(key, s.pool)
	if ok {
		atomic.AddInt64(&s.count, -1)
	}
}

func (s *atomicState[K, V]) pop(key K) (V, bool) {
	b := s.bucket(key)
	v, ok := b.delete(key, s.pool)
	if ok {
		atomic.AddInt64(&s.count, -1)
	}
	return v, ok
}

func (s *atomicState[K, V]) set(key K, value V) {
	b := s.bucket(key)
	ok := b.set(key, value, s.pool)
	if ok {
		atomic.AddInt64(&s.count, 1)
	}
}

func (s *atomicState[K, V]) swap(key K, value V) (v V, _ bool) {
	b := s.bucket(key)
	v, ok := b.swap(key, value, s.pool)
	if !ok {
		atomic.AddInt64(&s.count, 1)
	}
	return v, ok
}

func (s *atomicState[K, V]) range_(fn func(K, V) bool) bool {
	for i := range s.buckets {
		ok := s.buckets[i].range_(fn, s.pool)
		if !ok {
			return false
		}
	}
	return true
}

func (s *atomicState[K, V]) rangeLocked(fn func(K, V) bool) {
	for i := range s.buckets {
		ok := s.buckets[i].rangeLocked(fn)
		if !ok {
			return
		}
	}
}

// private

func (s *atomicState[K, V]) bucket(key K) *atomicBucket[K, V] {
	h := hashing.Hash(key)
	i := int(h) % len(s.buckets)
	return &s.buckets[i]
}
