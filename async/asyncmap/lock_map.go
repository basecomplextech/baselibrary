// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asyncmap

import (
	"runtime"
	"unsafe"

	"github.com/basecomplextech/baselibrary/async"
	"github.com/basecomplextech/baselibrary/internal/hashing"
	"github.com/basecomplextech/baselibrary/pools"
	"github.com/basecomplextech/baselibrary/status"
)

// LockMap holds locks for different keys.
//
// The map is a sharded map, which uses a lock per shard.
// The number of shards is equal to the number of CPU cores.
type LockMap[K comparable] interface {
	// Get returns a key lock, the lock must be freed after use.
	//
	// Usage:
	//
	//	m := NewLockMap[int]()
	//
	//	lock := m.Get(123)
	//	defer lock.Free()
	//
	//	select {
	//	case <-lock.Lock():
	//	case <-time.After(time.Second):
	//		return status.Timeout
	//	case <-ctx.Wait():
	//		return ctx.Status()
	//	}
	//	defer lock.Unlock()
	Get(key K) KeyLock

	// Contains returns true if the key is present.
	//
	// Usually it means that the key is locked, but it is not guaranteed.
	// In the latter case the key is unlocked but is not yet freed.
	Contains(key K) bool

	// Lock returns a locked key, the key must be freed after use.
	//
	// Usage:
	//
	//	m := NewLockMap[int]()
	//
	//	lock, st := m.Lock(ctx, 123)
	//	if !st.OK() {
	//		return st
	//	}
	//	defer lock.Free()
	Lock(ctx async.Context, key K) (LockedKey, status.Status)

	// LockMap locks the map itself, internally it locks all shards.
	//
	// Usage:
	//
	//	m := NewLockMap[int]()
	//
	//	locks := m.LockMap()
	//	defer locks.Free()
	//
	//	for key := range keys {
	//		ok := locks.Contains(key)
	//		// ...
	//	}
	LockMap() LockedMap[K]
}

// LockedMap is an interface to interact with a map locked in exclusive mode.
type LockedMap[K comparable] interface {
	// Contains returns true if the key is present.
	Contains(key K) bool

	// Range ranges over all keys.
	Range(f func(key K) bool)

	// Free unlocks the map itself, internally it unlocks all shards.
	Free()
}

// NewLockMap returns a new lock map.
func NewLockMap[K comparable]() LockMap[K] {
	return newLockMap[K]()
}

// internal

var _ LockMap[any] = (*lockMap[any])(nil)

type lockMap[K comparable] struct {
	shards []lockShard[K]
}

func newLockMap[K comparable]() *lockMap[K] {
	cpus := runtime.NumCPU()
	cpuLines := 16
	lineSize := 256

	size := unsafe.Sizeof(lockShard[K]{})
	total := cpus * cpuLines * lineSize
	n := int(total / int(size))

	pool := pools.NewPoolFunc[*lockItem[K]](newLockItem)
	shards := make([]lockShard[K], n)
	for i := range shards {
		shards[i] = newLockShard(pool)
	}

	return &lockMap[K]{
		shards: shards,
	}
}

// Get returns a key key, the lock must be freed after use.
func (m *lockMap[K]) Get(key K) KeyLock {
	// Get lock item
	shard := m.shard(key)
	item := shard.get(key)

	// Return key lock
	return &keyLock[K]{item}
}

// Contains returns true if the key is present.
//
// Usually it means that the key is locked, but it is not guaranteed.
// In the latter case the key is unlocked but is not yet freed.
func (m *lockMap[K]) Contains(key K) bool {
	shard := m.shard(key)
	return shard.contains(key)
}

// Lock returns a locked key, the key must be freed after use.
func (m *lockMap[K]) Lock(ctx async.Context, key K) (LockedKey, status.Status) {
	// Get lock item
	shard := m.shard(key)
	item := shard.get(key)

	// Free if not locked
	done := false
	defer func() {
		if !done {
			item.free()
		}
	}()

	// Try lock
	select {
	case <-item.lock:
	default:
		// Lock or wait
		// Context channel is lazily allocated, so try to postpone calling wait.
		select {
		case <-item.lock:
		case <-ctx.Wait():
			return nil, ctx.Status()
		}
	}

	// Return locked key
	k := &lockedKey[K]{item}
	done = true
	return k, status.OK
}

// LockMap locks the map itself, internally it locks all shards.
func (m *lockMap[K]) LockMap() LockedMap[K] {
	for i := range m.shards {
		shard := &m.shards[i]
		shard.mu.Lock()
	}

	return &lockedMap[K]{m}
}

// unlockMap unlocks the map itself, internally it unlocks all shards.
func (m *lockMap[K]) unlockMap() {
	for i := range m.shards {
		shard := &m.shards[i]
		shard.mu.Unlock()
	}
}

func (m *lockMap[K]) shard(key K) *lockShard[K] {
	index := hashing.Shard(key, len(m.shards))
	return &m.shards[index]
}

// locked

var _ LockedMap[any] = (*lockedMap[any])(nil)

type lockedMap[K comparable] struct {
	m *lockMap[K]
}

// Contains returns true if the key is present.
func (m *lockedMap[K]) Contains(key K) bool {
	shard := m.m.shard(key)
	return shard.containsLocked(key)
}

// Range ranges over all keys.
func (m *lockedMap[K]) Range(f func(key K) bool) {
	for i := range m.m.shards {
		shard := &m.m.shards[i]
		shard.rangeLocked(f)
	}
}

// Free unlocks the map itself, internally it unlocks all shards.
func (m *lockedMap[K]) Free() {
	m.m.unlockMap()
}
