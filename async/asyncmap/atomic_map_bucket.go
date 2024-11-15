// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package asyncmap

import (
	"sync"
	"sync/atomic"

	"github.com/basecomplextech/baselibrary/pools"
)

type atomicMapBucket[K comparable, V any] struct {
	wmu   sync.Mutex                           // write mutex
	ref   atomic.Int64                         // entry id and external refs packed into int64
	entry atomic.Pointer[atomicMapEntry[K, V]] // linked list of entries
}

func newAtomicMapBucket[K comparable, V any]() *atomicMapBucket[K, V] {
	return &atomicMapBucket[K, V]{}
}

func (b *atomicMapBucket[K, V]) get(key K, pool pools.Pool[*atomicMapEntry[K, V]]) (v V, _ bool) {
	// Acquire entry, increment external refs
	entry, prev, refs, ok := b.acquireEntry()
	if !ok {
		return
	}

	// Get item value
	v, ok = entry.get(key)

	// Release entry, decrement internal refs
	b.releaseEntry(entry, prev, pool)

	// Maybe reduce refs to avoid overflow
	b.reduceRefs(refs)
	return v, ok
}

func (b *atomicMapBucket[K, V]) getOrSet(key K, value V, pool pools.Pool[*atomicMapEntry[K, V]]) (
	v V, _ bool) {

	b.wmu.Lock()
	defer b.wmu.Unlock()

	// Load current entry
	entry := b.entry.Load()

	// Try to get value
	if entry != nil {
		v, ok := entry.get(key)
		if ok {
			return v, true
		}
	}

	// Make next entry
	next := newAtomicMapEntry(pool)
	next.init(entry)
	next.set(key, value)

	// Swap entry
	b.swapEntry(next, entry, pool)
	return value, false
}

func (b *atomicMapBucket[K, V]) set(key K, value V, pool pools.Pool[*atomicMapEntry[K, V]]) bool {
	b.wmu.Lock()
	defer b.wmu.Unlock()

	// Load current entry
	entry := b.entry.Load()

	// Make next entry
	next := newAtomicMapEntry(pool)
	next.init(entry)
	ok := next.set(key, value)

	// Swap entry
	b.swapEntry(next, entry, pool)
	return ok
}

func (b *atomicMapBucket[K, V]) setAbsent(key K, value V, pool pools.Pool[*atomicMapEntry[K, V]]) bool {
	b.wmu.Lock()
	defer b.wmu.Unlock()

	// Load current entry
	entry := b.entry.Load()

	// Check if exists
	if entry != nil {
		if _, ok := entry.get(key); ok {
			return false
		}
	}

	// Make next entry
	next := newAtomicMapEntry(pool)
	next.init(entry)
	next.set(key, value)

	// Swap entry
	b.swapEntry(next, entry, pool)
	return true
}

func (b *atomicMapBucket[K, V]) swap(key K, value V, pool pools.Pool[*atomicMapEntry[K, V]]) (
	v V, ok bool) {

	b.wmu.Lock()
	defer b.wmu.Unlock()

	// Load current entry
	entry := b.entry.Load()

	// Try to get value
	if entry != nil {
		v, ok = entry.get(key)
	}

	// Make next entry
	next := newAtomicMapEntry(pool)
	next.init(entry)
	next.set(key, value)

	// Swap entry
	b.swapEntry(next, entry, pool)
	return v, ok
}

func (b *atomicMapBucket[K, V]) delete(key K, pool pools.Pool[*atomicMapEntry[K, V]]) (v V, ok bool) {
	b.wmu.Lock()
	defer b.wmu.Unlock()

	// Load current entry
	entry := b.entry.Load()
	if entry == nil {
		return v, false
	}

	// Make next entry
	next := newAtomicMapEntry(pool)
	next.init(entry)
	v, ok = next.delete(key)

	// Swap entry
	b.swapEntry(next, entry, pool)
	return v, ok
}

func (b *atomicMapBucket[K, V]) range_(fn func(K, V) bool, pool pools.Pool[*atomicMapEntry[K, V]]) (
	continue_ bool) {

	// Acquire entry, increment external refs
	entry, prev, refs, ok := b.acquireEntry()
	if !ok {
		return true
	}

	// Iterate over entry items
	continue_ = entry.range_(fn)

	// Release entry, decrement internal refs
	b.releaseEntry(entry, prev, pool)

	// Maybe reduce refs to avoid overflow
	b.reduceRefs(refs)
	return continue_
}

func (b *atomicMapBucket[K, V]) rangeLocked(fn func(K, V) bool) (continue_ bool) {
	entry := b.entry.Load()
	if entry == nil {
		return true
	}
	return entry.range_(fn)
}

// private

func (b *atomicMapBucket[K, V]) acquireEntry() (
	entry *atomicMapEntry[K, V],
	prev *atomicMapEntry[K, V],
	refs int32,
	ok bool,
) {
	// Increment external refs
	ref := b.ref.Add(1)
	id, refs := unpackAtomicMapEntryRef(ref)
	if id == 0 {
		return nil, nil, 0, false
	}

	// Load current entry
	entry = b.entry.Load()
	if entry == nil {
		return nil, nil, 0, false
	}

	// Find entry in linked list
	for entry.id != id {
		prev = entry

		entry = entry.prev.Load()
		if entry == nil {
			return nil, nil, 0, false
		}
	}

	return entry, prev, refs, true
}

func (b *atomicMapBucket[K, V]) releaseEntry(
	entry *atomicMapEntry[K, V],
	prev *atomicMapEntry[K, V],
	pool pools.Pool[*atomicMapEntry[K, V]],
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

func (b *atomicMapBucket[K, V]) swapEntry(
	next *atomicMapEntry[K, V],
	prev *atomicMapEntry[K, V],
	pool pools.Pool[*atomicMapEntry[K, V]],
) {
	// Store entry
	b.entry.Store(next)

	// Swap reference
	nextRef := packAtomicMapEntryRef(next.id, 1)
	lastRef := b.ref.Swap(nextRef)

	// Return if no previous
	if prev == nil {
		return
	}

	// Increment previous internal refs by (external-1)
	_, ext := unpackAtomicMapEntryRef(lastRef)
	int_ := prev.refs.Add(ext - 1)
	if int_ != 0 {
		return
	}

	// Free previous when refs are 0
	next.prev.Store(prev.prev.Load())

	prev.reset()
	pool.Put(prev)
}

// reduceRefs reduces entry refs when they reach 1000_000_000 to avoid overflow.
func (b *atomicMapBucket[K, V]) reduceRefs(refs int32) {
	if refs != 1000_000_000 {
		return
	}

	b.wmu.Lock()
	defer b.wmu.Unlock()

	// Load external refs
	ref := b.ref.Load()
	_, ext := unpackAtomicMapEntryRef(ref)
	if ext < 1000_000_000 {
		return
	}

	// Decrement external refs
	b.ref.Add(-999_000_000)

	// Load current entry
	entry := b.entry.Load()
	if entry == nil {
		return
	}

	// Increment internal refs
	entry.refs.Add(999_000_000)
}
