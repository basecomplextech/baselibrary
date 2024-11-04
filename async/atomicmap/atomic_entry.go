// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package atomics

import (
	"sync/atomic"

	"github.com/basecomplextech/baselibrary/collect/slices2"
)

type atomicEntry[K comparable, V any] struct {
	id   int32
	refs atomic.Int32 // internal refcount, 0 by default, external is 1
	prev atomic.Pointer[atomicEntry[K, V]]

	item atomicItem[K, V]
	more []atomicItem[K, V]
}

type atomicItem[K comparable, V any] struct {
	set   bool
	key   K
	value V
}

// init inits a new entry, copies the previous items if any.
func (e *atomicEntry[K, V]) init(prev *atomicEntry[K, V]) {
	e.id = 1
	if prev == nil {
		return
	}

	e.id = prev.id + 1
	e.prev.Store(prev)
	e.item = prev.item

	if len(prev.more) > 0 {
		e.more = append(e.more, prev.more...)
	}
}

// items

func (e *atomicEntry[K, V]) get(key K) (v V, ok bool) {
	if e.item.set {
		if e.item.key == key {
			return e.item.value, true
		}
	}

	for i := range e.more {
		item := &e.more[i]
		if item.key == key {
			return item.value, true
		}
	}

	return v, false
}

func (e *atomicEntry[K, V]) set(key K, value V) bool {
	if !e.item.set {
		e.item = atomicItem[K, V]{true, key, value}
		return true
	}

	if e.item.key == key {
		e.item.value = value
		return false
	}

	for i := range e.more {
		item := &e.more[i]
		if item.key == key {
			item.value = value
			return false
		}
	}

	item := atomicItem[K, V]{true, key, value}
	e.more = append(e.more, item)
	return true
}

func (e *atomicEntry[K, V]) delete(key K) (v V, ok bool) {
	if e.item.set && e.item.key == key {
		v = e.item.value
		e.item = atomicItem[K, V]{}
		return v, true
	}

	for i := range e.more {
		item := &e.more[i]
		if item.key == key {
			v = item.value
			e.more = slices2.RemoveAt(e.more, i, 1)
			return v, true
		}
	}

	return v, false
}

func (e *atomicEntry[K, V]) range_(fn func(K, V) bool) (continue_ bool) {
	if e.item.set {
		if !fn(e.item.key, e.item.value) {
			return false
		}
	}

	for i := range e.more {
		item := &e.more[i]
		if !fn(item.key, item.value) {
			return false
		}
	}

	return true
}

// reset

func (e *atomicEntry[K, V]) reset() {
	more := e.more

	*e = atomicEntry[K, V]{}
	e.more = slices2.Truncate(more)
}

// ref

func packAtomicEntryRef(id int32, refcount int32) int64 {
	return int64(id)<<32 | int64(refcount)
}

func unpackAtomicEntryRef(r int64) (id int32, refcount int32) {
	return int32(r >> 32), int32(r & 0xffffffff)
}
