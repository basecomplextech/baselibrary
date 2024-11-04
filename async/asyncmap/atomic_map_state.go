// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asyncmap

import (
	"sync/atomic"

	"github.com/basecomplextech/baselibrary/pools"
)

type atomicMapState[K comparable, V any] struct {
	pool pools.Pool[*atomicMapEntry[K, V]]

	count     int64 // number of items
	threshold int   // next resize threshold

	buckets []atomicMapBucket[K, V]
}

func newAtomicMapState[K comparable, V any](
	size int,
	pool pools.Pool[*atomicMapEntry[K, V]],
) *atomicMapState[K, V] {

	size = max(size, atomicMapMinSize)
	threshold := int(float64(size) * atomicMapThreshold)

	return &atomicMapState[K, V]{
		pool:      pool,
		threshold: threshold,

		buckets: make([]atomicMapBucket[K, V], size),
	}
}

func emptyAtomicMapState[K comparable, V any](s *atomicMapState[K, V]) *atomicMapState[K, V] {
	return &atomicMapState[K, V]{
		pool:      s.pool,
		threshold: s.threshold,

		buckets: make([]atomicMapBucket[K, V], len(s.buckets)),
	}
}

// methods

func (s *atomicMapState[K, V]) contains(h uint32, key K) bool {
	b := s.bucket(h)
	_, ok := b.get(key, s.pool)
	return ok
}

func (s *atomicMapState[K, V]) get(h uint32, key K) (V, bool) {
	b := s.bucket(h)
	return b.get(key, s.pool)
}

func (s *atomicMapState[K, V]) getOrSet(h uint32, key K, value V) (V, bool) {
	b := s.bucket(h)
	v, ok := b.getOrSet(key, value, s.pool)
	if !ok {
		atomic.AddInt64(&s.count, 1)
	}
	return v, ok
}

func (s *atomicMapState[K, V]) delete(h uint32, key K) {
	b := s.bucket(h)
	_, ok := b.delete(key, s.pool)
	if ok {
		atomic.AddInt64(&s.count, -1)
	}
}

func (s *atomicMapState[K, V]) pop(h uint32, key K) (V, bool) {
	b := s.bucket(h)
	v, ok := b.delete(key, s.pool)
	if ok {
		atomic.AddInt64(&s.count, -1)
	}
	return v, ok
}

func (s *atomicMapState[K, V]) set(h uint32, key K, value V) {
	b := s.bucket(h)
	ok := b.set(key, value, s.pool)
	if ok {
		atomic.AddInt64(&s.count, 1)
	}
}

func (s *atomicMapState[K, V]) swap(h uint32, key K, value V) (v V, _ bool) {
	b := s.bucket(h)
	v, ok := b.swap(key, value, s.pool)
	if !ok {
		atomic.AddInt64(&s.count, 1)
	}
	return v, ok
}

func (s *atomicMapState[K, V]) range_(fn func(K, V) bool) bool {
	for i := range s.buckets {
		ok := s.buckets[i].range_(fn, s.pool)
		if !ok {
			return false
		}
	}
	return true
}

func (s *atomicMapState[K, V]) rangeLocked(fn func(K, V) bool) {
	for i := range s.buckets {
		ok := s.buckets[i].rangeLocked(fn)
		if !ok {
			return
		}
	}
}

// private

func (s *atomicMapState[K, V]) bucket(h uint32) *atomicMapBucket[K, V] {
	i := int(h) % len(s.buckets)
	return &s.buckets[i]
}
