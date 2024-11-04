// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package atomics

import (
	"sync"
	"sync/atomic"

	"github.com/basecomplextech/baselibrary/pools"
)

type atomicBucket[K comparable, V any] struct {
	wmu   sync.Mutex                        // write mutex
	ref   atomic.Int64                      // entry id and external refs packed into int64
	entry atomic.Pointer[atomicEntry[K, V]] // linked list of entries
}

func newAtomicBucket[K comparable, V any]() *atomicBucket[K, V] {
	return &atomicBucket[K, V]{}
}

func (b *atomicBucket[K, V]) get(key K, pool pools.Pool[*atomicEntry[K, V]]) (v V, _ bool) {
	// Acquire entry, increment external refs
	entry, prev, ok := b.acquireEntry()
	if !ok {
		return
	}

	// Get item value
	v, ok = entry.get(key)

	// Release entry, decrement internal refs
	b.releaseEntry(entry, prev, pool)
	return v, ok
}

func (b *atomicBucket[K, V]) getOrSet(key K, value V, pool pools.Pool[*atomicEntry[K, V]]) (v V, _ bool) {
	b.wmu.Lock()
	defer b.wmu.Unlock()

	// Load current entry
	entry := b.entry.Load()

	// Try to get value
	v, ok := entry.get(key)
	if ok {
		return v, true
	}

	// Make next entry
	next := pool.New()
	next.init(entry)
	next.set(key, value)

	// Swap entry
	b.swapEntry(next, entry, pool)
	return v, false
}

func (b *atomicBucket[K, V]) set(key K, value V, pool pools.Pool[*atomicEntry[K, V]]) bool {
	b.wmu.Lock()
	defer b.wmu.Unlock()

	// Load current entry
	entry := b.entry.Load()

	// Make next entry
	next := pool.New()
	next.init(entry)
	ok := next.set(key, value)

	// Swap entry
	b.swapEntry(next, entry, pool)
	return ok
}

func (b *atomicBucket[K, V]) swap(key K, value V, pool pools.Pool[*atomicEntry[K, V]]) (v V, ok bool) {
	b.wmu.Lock()
	defer b.wmu.Unlock()

	// Load current entry
	entry := b.entry.Load()
	v, ok = entry.get(key)

	// Make next entry
	next := pool.New()
	next.init(entry)
	next.set(key, value)

	// Swap entry
	b.swapEntry(next, entry, pool)
	return v, ok
}

func (b *atomicBucket[K, V]) delete(key K, pool pools.Pool[*atomicEntry[K, V]]) (v V, ok bool) {
	b.wmu.Lock()
	defer b.wmu.Unlock()

	// Load current entry
	entry := b.entry.Load()

	// Make next entry
	next := pool.New()
	next.init(entry)
	v, ok = next.delete(key)

	// Swap entry
	b.swapEntry(next, entry, pool)
	return v, ok
}

func (b *atomicBucket[K, V]) range_(fn func(K, V) bool, pool pools.Pool[*atomicEntry[K, V]]) (
	continue_ bool) {

	// Acquire entry, increment external refs
	entry, prev, ok := b.acquireEntry()
	if !ok {
		return true
	}

	// Iterate over entry items
	continue_ = entry.range_(fn)

	// Release entry, decrement internal refs
	b.releaseEntry(entry, prev, pool)
	return continue_
}

func (b *atomicBucket[K, V]) rangeLocked(fn func(K, V) bool) (continue_ bool) {
	entry := b.entry.Load()
	if entry == nil {
		return true
	}

	return entry.range_(fn)
}

// private

func (b *atomicBucket[K, V]) acquireEntry() (
	entry *atomicEntry[K, V],
	prev *atomicEntry[K, V],
	ok bool,
) {
	// Increment external refs
	ref := b.ref.Add(1)
	id, _ := unpackAtomicEntryRef(ref)
	if id == 0 {
		return nil, nil, false
	}

	// Load current entry
	entry = b.entry.Load()
	if entry == nil {
		return nil, nil, false
	}

	// Find entry in linked list
	for entry.id != id {
		prev = entry

		entry = entry.prev.Load()
		if entry == nil {
			return nil, nil, false
		}
	}

	return entry, prev, true
}

func (b *atomicBucket[K, V]) releaseEntry(
	entry *atomicEntry[K, V],
	prev *atomicEntry[K, V],
	pool pools.Pool[*atomicEntry[K, V]],
) {
	// Decrement internal refs
	int_ := entry.refs.Add(-1)
	if int_ != 0 {
		return
	}

	// Free entry when refs are 0
	if prev != nil {
		prev.prev.Store(entry.prev.Load())
	}

	entry.reset()
	pool.Put(entry)
}

func (b *atomicBucket[K, V]) swapEntry(
	next *atomicEntry[K, V],
	prev *atomicEntry[K, V],
	pool pools.Pool[*atomicEntry[K, V]],
) {
	// Store entry
	b.entry.Store(next)

	// Swap reference
	nextRef := packAtomicEntryRef(next.id, 1)
	lastRef := b.ref.Swap(nextRef)

	// Return if no previous
	if prev == nil {
		return
	}

	// Increment previous internal refs by (external-1)
	_, ext := unpackAtomicEntryRef(lastRef)
	int_ := prev.refs.Add(ext - 1)
	if int_ != 0 {
		return
	}

	// Free previous when refs are 0
	next.prev.Store(prev.prev.Load())

	prev.reset()
	pool.Put(prev)
}
