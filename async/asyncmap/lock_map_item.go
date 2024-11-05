// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asyncmap

import "github.com/basecomplextech/baselibrary/pools"

type lockMapItem[K comparable] struct {
	b *lockMapBucket[K]

	refs int32
	lock chan struct{}

	key K
}

func newLockMapItem[K comparable](b *lockMapBucket[K], key K) *lockMapItem[K] {
	m, ok := b.m.pool.Get()
	if !ok {
		m = makeLockMapItem[K]()
	}

	m.b = b
	m.key = key
	m.refs = 1
	return m
}

func makeLockMapItem[K comparable]() *lockMapItem[K] {
	m := &lockMapItem[K]{}
	m.lock = make(chan struct{}, 1)
	m.lock <- struct{}{}
	return m
}

func (m *lockMapItem[K]) unlock() {
	select {
	case m.lock <- struct{}{}:
	default:
		panic("unlock of unlocked key lock")
	}
}

func (m *lockMapItem[K]) release() {
	deleted := m.b.release(m)
	if !deleted {
		return
	}

	b := m.b
	m.reset()
	b.m.pool.Put(m)
}

// private

func (m *lockMapItem[K]) reset() {
	lock := m.lock
	select {
	case m.lock <- struct{}{}:
	default:
	}

	*m = lockMapItem[K]{}
	m.lock = lock
}

// pools

var lockMapItemPools = pools.NewPools()

func newLockMapPool[K comparable]() pools.Pool[*lockMapItem[K]] {
	return pools.GetPool[*lockMapItem[K]](lockMapItemPools)
}
