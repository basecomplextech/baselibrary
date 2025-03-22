// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asyncmap

import (
	"runtime"
	"unsafe"

	"github.com/basecomplextech/baselibrary/async"
	"github.com/basecomplextech/baselibrary/hashing"
	"github.com/basecomplextech/baselibrary/pools"
	"github.com/basecomplextech/baselibrary/status"
)

// LockMap holds locks for different keys.
//
// The map uses multiple buckets (shards) each with its own mutex.
// Buckets are stored in multiple cache lines to try to avoid false sharing.
// The number of cache lines is equal to the number of CPUs.
//
// # Benchmarks
//
//	cpu: Apple M1 Pro
//	BenchmarkLockMap_Get-10              	20442770	        53.08 ns/op	        18.84 mops	       8 B/op	       1 allocs/op
//	BenchmarkLockMap_Get_Parallel-10     	18444724	        66.06 ns/op	        15.14 mops	       8 B/op	       1 allocs/op
//	BenchmarkLockMap_Lock-10             	12843542	        93.09 ns/op	        10.74 mops	       8 B/op	       1 allocs/op
//	BenchmarkLockMap_Lock_Parallel-10    	14246732	        88.46 ns/op	        11.30 mops	       8 B/op	       1 allocs/op
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

	// LockMap locks the map itself, internally it locks all buckets.
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
	LockMap() LockedLockMap[K]
}

// NewLockMap returns a new lock map.
func NewLockMap[K comparable]() LockMap[K] {
	return newLockMap[K]()
}

// internal

var _ LockMap[any] = (*lockMap[any])(nil)

type lockMap[K comparable] struct {
	pool    pools.Pool[*lockMapItem[K]]
	buckets []lockMapBucket[K]
}

func newLockMap[K comparable]() *lockMap[K] {
	pool := newLockMapPool[K]()

	// CPUs
	cpus := runtime.NumCPU()
	cacheLine := 256

	// Calculate number of buckets
	bucketSize := unsafe.Sizeof(lockMapBucket[K]{})
	bucketNum := (cpus * cacheLine) / int(bucketSize)

	// Make map
	m := &lockMap[K]{
		pool:    pool,
		buckets: make([]lockMapBucket[K], bucketNum),
	}

	// Init buckets
	for i := range m.buckets {
		m.buckets[i].init(m)
	}
	return m
}

// Get returns a key key, the lock must be freed after use.
func (m *lockMap[K]) Get(key K) KeyLock {
	// Get lock item
	b := m.bucket(key)
	item := b.get(key)

	// Return key lock
	return newLockMapKeyLock(item)
}

// Contains returns true if the key is present.
//
// Usually it means that the key is locked, but it is not guaranteed.
// In the latter case the key is unlocked but is not yet freed.
func (m *lockMap[K]) Contains(key K) bool {
	b := m.bucket(key)
	return b.contains(key)
}

// Lock returns a locked key, the key must be freed after use.
func (m *lockMap[K]) Lock(ctx async.Context, key K) (LockedKey, status.Status) {
	// Get lock item
	b := m.bucket(key)
	item := b.get(key)

	// Release if not locked
	done := false
	defer func() {
		if !done {
			item.release()
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
	k := newLockMapLockedKey(item)
	done = true
	return k, status.OK
}

// LockMap locks the map itself, internally it locks all buckets.
func (m *lockMap[K]) LockMap() LockedLockMap[K] {
	i := 0
	done := false
	defer func() {
		if !done {
			for ; i >= 0; i-- {
				m.buckets[i].mu.Unlock()
			}
		}
	}()

	for i = range m.buckets {
		m.buckets[i].mu.Lock()
	}

	locked := newLockedLockMap(m)
	done = true
	return locked
}

// internal

// unlockMap unlocks the map itself, internally it unlocks all buckets.
func (m *lockMap[K]) unlockMap() {
	for i := range m.buckets {
		b := &m.buckets[i]
		b.mu.Unlock()
	}
}

func (m *lockMap[K]) bucket(key K) *lockMapBucket[K] {
	h := hashing.Hash(key)
	i := int(h) % len(m.buckets)
	return &m.buckets[i]
}
