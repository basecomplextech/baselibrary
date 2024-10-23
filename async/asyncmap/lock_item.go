// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asyncmap

type lockItem[K comparable] struct {
	shard *lockShard[K]
	lock  chan struct{}
	refs  int32

	key K
}

func newLockItem[K comparable]() *lockItem[K] {
	m := &lockItem[K]{
		lock: make(chan struct{}, 1),
		refs: 1,
	}
	m.lock <- struct{}{}
	return m
}

func (m *lockItem[K]) unlock() {
	select {
	case m.lock <- struct{}{}:
	default:
		panic("unlock of unlocked key lock")
	}
}

func (m *lockItem[K]) free() {
	deleted := m.shard.free(m)
	if deleted {
		m.release()
	}
}

func (m *lockItem[K]) release() {
	s := m.shard
	m.reset()
	s.pool.Put(m)
}

func (m *lockItem[K]) reset() {
	var zero K
	m.shard = nil
	m.key = zero
	m.refs = 0

	select {
	case m.lock <- struct{}{}:
	default:
	}
}
