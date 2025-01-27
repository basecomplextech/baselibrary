// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package immap

import (
	"github.com/basecomplextech/baselibrary/compare"
	"github.com/basecomplextech/baselibrary/pools"
	"github.com/basecomplextech/baselibrary/ref"
)

// Map is an immutable sorted map, implemented as a btree.
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

	// Get returns an item by a key.
	Get(key K) (V, bool)

	// Contains returns true if a key exists.
	Contains(key K) bool

	// Iterator returns an iterator.
	Iterator() Iterator[K, V]

	// Keys returns all keys.
	Keys() []K

	// Write

	// Set adds an item to the map.
	Set(key K, value V)

	// Move moves an item from one key to another, or returns false if the key does not exist.
	Move(key K, newKey K) bool

	// Delete deletes an item by a key.
	Delete(key K)

	// Internal

	// Free frees the map.
	Free()
}

// Item is a map item.
type Item[K, V any] struct {
	Key   K
	Value V
}

// New returns an empty map.
func New[K, V any](mutable bool, compare compare.Compare[K]) Map[K, V] {
	m := newBtree[K, V](compare)
	if !mutable {
		m.Freeze()
	}
	return m
}

// New returns an empty map wrapped in a ref.
func NewRef[K, V any](mutable bool, compare compare.Compare[K]) ref.R[Map[K, V]] {
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
	compare compare.Compare[K]

	root    node[K, V]
	mod     int // track concurrent modifications
	height  int
	length  int64
	mutable bool
}

func newBtree[K, V any](compare compare.Compare[K]) *btree[K, V] {
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
		panic("cannot clone mutable immap")
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

// Get returns an item by a key.
func (t *btree[K, V]) Get(key K) (v V, ok bool) {
	return t.root.get(key, t.compare)
}

// Contains returns true if a key exists.
func (t *btree[K, V]) Contains(key K) bool {
	return t.root.contains(key, t.compare)
}

// Iterator returns an iterator.
func (t *btree[K, V]) Iterator() Iterator[K, V] {
	it := newIterator(t)
	it.SeekToStart()
	return it
}

// Keys returns all keys.
func (t *btree[K, V]) Keys() []K {
	n := t.length
	if n == 0 {
		return nil
	}

	it := t.Iterator()
	defer it.Free()

	keys := make([]K, 0, n)
	for it.Next() {
		key := it.Key()
		keys = append(keys, key)
	}
	return keys
}

// Write

// Set adds an item to the map.
func (t *btree[K, V]) Set(key K, value V) {
	if !t.mutable {
		panic("operation on immutable immap")
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

// Move moves an item from one key to another, or returns false if the key does not exist.
func (t *btree[K, V]) Move(key K, newKey K) bool {
	if !t.mutable {
		panic("operation on immutable immap")
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
	return true
}

// Delete deletes an item by a key.
func (t *btree[K, V]) Delete(key K) {
	if !t.mutable {
		panic("operation on immutable immap")
	}

	// Mutate root
	node := t.mutateRoot()

	// Delete item
	_, mod := node.delete(key, t.compare)
	if !mod {
		return
	}

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
		panic("operation on immutable map")
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
}

// for tests

// items returns items as a slice.
func (t *btree[K, V]) items() []Item[K, V] {
	result := make([]Item[K, V], 0, t.length)

	// LIFO stack
	stack := []node[K, V]{t.root}
	for len(stack) > 0 {
		// Pop node
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		switch n := node.(type) {
		case *branchNode[K, V]:
			// Push in reverse order
			for i := len(n.items) - 1; i >= 0; i-- {
				item := n.items[i]
				stack = append(stack, item.node)
			}

		case *leafNode[K, V]:
			for _, item := range n.items {
				item1 := Item[K, V]{
					Key:   item.key,
					Value: item.value,
				}
				result = append(result, item1)
			}
		}
	}
	return result
}

// keys returns item keys as a slice.
func (t *btree[K, V]) keys() []K {
	result := make([]K, 0, t.length)

	// LIFO stack
	stack := []node[K, V]{t.root}
	for len(stack) > 0 {
		// Pop node
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		switch n := node.(type) {
		case *branchNode[K, V]:
			// Push in reverse order
			for i := len(n.items) - 1; i >= 0; i-- {
				item := n.items[i]
				stack = append(stack, item.node)
			}

		case *leafNode[K, V]:
			for _, item := range n.items {
				result = append(result, item.key)
			}
		}
	}
	return result
}

// values returns item values as a slice.
func (t *btree[K, V]) values() []V {
	result := make([]V, 0, t.length)

	// LIFO stack
	stack := []node[K, V]{t.root}
	for len(stack) > 0 {
		// Pop node
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		switch n := node.(type) {
		case *branchNode[K, V]:
			// Push in reverse order
			for i := len(n.items) - 1; i >= 0; i-- {
				item := n.items[i]
				stack = append(stack, item.node)
			}

		case *leafNode[K, V]:
			for _, item := range n.items {
				result = append(result, item.value)
			}
		}
	}
	return result
}

// state pool

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
