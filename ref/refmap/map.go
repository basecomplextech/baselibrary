// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package refmap

import (
	"github.com/basecomplextech/baselibrary/iterator"
	"github.com/basecomplextech/baselibrary/pools"
	"github.com/basecomplextech/baselibrary/ref"
)

// Map is an immutable sorted map which stores countable references, implemented as a btree.
//
// The map retains/releases references internally, but does not retain them when iterating
// or returning.
type Map[K, V any] interface {
	// Empty returns true if the map is empty.
	Empty() bool

	// Length returns the number of items in this map, this can be an estimate.
	Length() int64

	// Mutable returns true if the map is mutable.
	Mutable() bool

	// Clone

	// Clone returns a mutable clone of the map.
	Clone() Map[K, V]

	// Freeze makes the map immutable.
	Freeze()

	// Read

	// Get returns an item by a key, does not retain the value.
	Get(key K) (ref.R[V], bool)

	// First returns the first item, does not retain the value.
	First() (K, ref.R[V], bool)

	// Last returns the last item, does not retain the value.
	Last() (K, ref.R[V], bool)

	// Contains returns true if the map contains a key.
	Contains(key K) bool

	// Iterator

	// Keys returns a key iterator.
	Keys() iterator.Iter[K]

	// Values returns a value iterator.
	Values() iterator.Iter[V]

	// Iterator returns a map iterator.
	Iterator() Iterator[K, V]

	// Write

	// SetNewRef adds an item to the map, wraps into into a reference.
	SetNewRef(key K, value V)

	// SetRetain adds an item reference to the map, and retains it.
	SetRetain(key K, value ref.R[V])

	// SetNoRetain adds an item reference to the map, does not retain it.
	SetNoRetain(key K, value ref.R[V])

	// Move moves an item from one key to another, or returns false if the key does not exist.
	Move(key K, newKey K) bool

	// Delete deletes an item by a key, releases its value.
	Delete(key K)

	// Internal

	// Free frees the map, releases all values.
	Free()
}

// CompareFunc compares two keys, and returns -1 if a < b, 0 if a == b, 1 if a > b.
type CompareFunc[K any] func(a, b K) int

// New returns an empty map.
func New[K, V any](mutable bool, compare CompareFunc[K]) Map[K, V] {
	m := newBtree[K, V](compare)
	if !mutable {
		m.Freeze()
	}
	return m
}

// New returns an empty map wrapped in a ref.
func NewRef[K, V any](mutable bool, compare CompareFunc[K]) ref.R[Map[K, V]] {
	m := New[K, V](mutable, compare)
	return ref.New(m)
}

// internal

const maxItems = 16

var _ Map[any, any] = (*btree[any, any])(nil)

type btree[K, V any] struct {
	*state[K, V]
}

type state[K, V any] struct {
	compare CompareFunc[K]

	root    node[K, V]
	mod     int // track concurrent modifications
	height  int
	length  int64
	mutable bool
}

func newBtree[K, V any](compare CompareFunc[K]) *btree[K, V] {
	b := &btree[K, V]{acquireState[K, V]()}
	b.compare = compare
	b.root = newLeafNode[K, V]()
	b.height = 1

	b.mutable = true
	return b
}

func (s *state[K, V]) reset() {
	*s = state[K, V]{}
}

// Empty returns true if the map is empty.
func (t *btree[K, V]) Empty() bool {
	return t.length == 0
}

// Length returns the number of items in this map, this can be an estimate.
func (t *btree[K, V]) Length() int64 {
	return t.length
}

// Mutable returns true if the map is mutable.
func (t *btree[K, V]) Mutable() bool {
	return t.mutable
}

// Clone

// Clone returns a mutable clone of the map.
func (t *btree[K, V]) Clone() Map[K, V] {
	if t.mutable {
		panic("cannot clone mutable refmap")
	}

	root1 := t.root.clone()
	t1 := &btree[K, V]{acquireState[K, V]()}
	t1.compare = t.compare

	t1.root = root1
	t1.height = t.height
	t1.length = t.length
	t1.mutable = true
	return Map[K, V](t1)
}

// Freeze makes the map immutable.
func (t *btree[K, V]) Freeze() {
	if !t.mutable {
		return
	}

	t.mutable = false
	t.root.freeze()
}

// Items

// Get returns an item by a key, does not retain the value.
func (t *btree[K, V]) Get(key K) (v ref.R[V], ok bool) {
	return t.root.get(key, t.compare)
}

// First returns the first item, does not retain the value.
func (t *btree[K, V]) First() (K, ref.R[V], bool) {
	return t.root.first()
}

// Last returns the last item, does not retain the value.
func (t *btree[K, V]) Last() (K, ref.R[V], bool) {
	return t.root.last()
}

// Contains returns true if the map contains a key.
func (t *btree[K, V]) Contains(key K) bool {
	return t.root.contains(key, t.compare)
}

// Iterator

// Keys returns a key iterator.
func (t *btree[K, V]) Keys() iterator.Iter[K] {
	it := newIterator(t)
	it.SeekToStart()
	return iterator.MapToKeys(it)
}

// Values returns a value iterator.
func (t *btree[K, V]) Values() iterator.Iter[V] {
	it := newIterator(t)
	it.SeekToStart()
	return iterator.MapToValues(it)
}

// Iterator returns a map iterator.
func (t *btree[K, V]) Iterator() Iterator[K, V] {
	it := newIterator(t)
	it.SeekToStart()
	return it
}

// Write

// SetNewRef adds an item to the map, wraps into into a reference.
func (t *btree[K, V]) SetNewRef(key K, value V) {
	var r ref.R[V]

	v, ok := (any)(value).(ref.Freer)
	if ok {
		r = ref.NewFreer(value, v)
	} else {
		r = ref.NewNoop(value)
	}

	t.SetRetain(key, r)
	r.Release()
}

// SetRetain adds an item reference to the map, and retains it.
func (t *btree[K, V]) SetRetain(key K, value ref.R[V]) {
	if !t.mutable {
		panic("operation on immutable refmap")
	}

	// Split if full
	t.maybeSplitRoot()

	// Mutate root
	node := t.mutateRoot()

	// Insert item
	mod := node.insert(key, value, t.compare)
	if !mod {
		return
	}

	// Increment length
	t.mod++
	t.length++
}

// SetNoRetain adds an item reference to the map, does not retain it.
func (t *btree[K, V]) SetNoRetain(key K, value ref.R[V]) {
	t.SetRetain(key, value)
	value.Release()
}

// Move moves an item from one key to another, or returns false if the key does not exist.
func (t *btree[K, V]) Move(key K, newKey K) bool {
	if !t.mutable {
		panic("operation on immutable refmap")
	}

	// Mutate root
	node := t.mutateRoot()

	// Delete item
	value, mod := node.delete(key, t.compare)
	if !mod {
		return false
	}

	// Insert item
	node.insert(newKey, value, t.compare)
	value.Release()
	return true
}

// Delete deletes an item by a key, releases its value.
func (t *btree[K, V]) Delete(key K) {
	if !t.mutable {
		panic("operation on immutable refmap")
	}

	// Mutate root
	node := t.mutateRoot()

	// Delete item
	value, mod := node.delete(key, t.compare)
	if !mod {
		return
	}
	value.Release()

	// Decrement length
	t.mod++
	t.length--
}

// Internal

// Free frees the map, releases all values.
func (t *btree[K, V]) Free() {
	t.root.release()
	t.root = nil

	t.height = 0
	t.length = 0

	// Release state
	s := t.state
	t.state = nil
	releaseState[K, V](s)
}

// root

func (t *btree[K, V]) mutateRoot() node[K, V] {
	if !t.mutable {
		panic("operation on immutable refmap")
	}

	if t.root.mutable() {
		return t.root
	}

	// Clone and replace root
	prev := t.root
	next := t.root.clone()
	t.root = next
	t.mod++

	// Release previous
	prev.release()
	return next
}

func (t *btree[K, V]) maybeSplitRoot() {
	if t.root.length() < maxItems {
		return
	}

	// Split root
	node := t.mutateRoot()
	next, ok := node.split()
	if !ok {
		return
	}

	// Make new root, move children to it
	t.root = newBranchNode(node, next)
	t.height++
	t.mod++
	return
}

// pools

var statePools = pools.NewPools()

func acquireState[K, V any]() *state[K, V] {
	v, ok := pools.Acquire[*state[K, V]](statePools)
	if ok {
		return v
	}
	return &state[K, V]{}
}

func releaseState[K, V any](s *state[K, V]) {
	s.reset()
	pools.Release(statePools, s)
}
