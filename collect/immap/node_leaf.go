// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package immap

import (
	"sort"
	"sync/atomic"

	"github.com/basecomplextech/baselibrary/collect/slices2"
	"github.com/basecomplextech/baselibrary/compare"
	"github.com/basecomplextech/baselibrary/pools"
)

var _ node[any, any] = (*leafNode[any, any])(nil)

type leafNode[K, V any] struct {
	items  []leafItem[K, V]
	_items [maxItems]leafItem[K, V]

	mut  bool
	refs int64
}

type leafItem[K, V any] struct {
	key   K
	value V
}

// newLeafNode returns a new mutable node.
func newLeafNode[K, V any](items ...Item[K, V]) *leafNode[K, V] {
	// Make node
	n := acquireLeaf[K, V]()
	n.items = n._items[:0]
	n.mut = true
	n.refs = 1

	// Copy items
	for _, item := range items {
		n.items = append(n.items, leafItem[K, V]{
			key:   item.Key,
			value: item.Value,
		})
	}
	return n
}

// cloneLeafNode returns a mutable node clone.
func cloneLeafNode[K, V any](n *leafNode[K, V]) *leafNode[K, V] {
	// Copy node
	n1 := acquireLeaf[K, V]()
	n1.items = n1._items[:0]
	n1.mut = true
	n1.refs = 1

	n1.items = n1.items[:len(n.items)]
	copy(n1.items, n.items)
	return n1
}

// nextLeafNode returns a new mutable node on a split, moves items to it.
func nextLeafNode[K, V any](items []leafItem[K, V]) *leafNode[K, V] {
	n := acquireLeaf[K, V]()
	n.items = n._items[:0]
	n.refs = 1
	n.mut = true

	n.items = n.items[:len(items)]
	copy(n.items, items)
	return n
}

// length returns the number of items in the node.
func (n *leafNode[K, V]) length() int {
	return len(n.items)
}

// minKey returns the minimum key in the node.
func (n *leafNode[K, V]) minKey() K {
	return n.items[0].key
}

// maxKey returns the maximum key in the node.
func (n *leafNode[K, V]) maxKey() K {
	return n.items[len(n.items)-1].key
}

// mutable returns true if the node is mutable.
func (n *leafNode[K, V]) mutable() bool {
	return n.mut
}

// get/insert/delete

// get returns for an item by key, or false if not found.
func (n *leafNode[K, V]) get(key K, compare compare.Compare[K]) (v V, ok bool) {
	index := n.indexOf(key, compare)

	// Return if not found
	if index >= len(n.items) {
		return
	}
	if compare(n.items[index].key, key) != 0 {
		return
	}

	item := n.items[index]
	return item.value, true
}

// insert inserts or updates an item, returns true if inserted.
func (n *leafNode[K, V]) insert(key K, value V, compare compare.Compare[K]) bool {
	if !n.mut {
		panic("operation on immutable node")
	}
	if len(n.items) == maxItems {
		panic("cannot insert into full node")
	}

	// Find item by key
	index := n.indexOf(key, compare)

	// Replace existing if found
	if index < len(n.items) {
		item := &n.items[index]

		// Swap item
		if compare(item.key, key) == 0 {
			item.value = value
			return false
		}
	}

	// Grow capacity
	if cap(n.items) < len(n.items)+1 {
		new := 2*len(n.items) + 1
		items := make([]leafItem[K, V], len(n.items), new)

		copy(items, n.items)
		n.items = items
	}

	// Shift greater items right
	n.items = n.items[:len(n.items)+1]
	copy(n.items[index+1:], n.items[index:])

	// Insert new item at index
	n.items[index] = leafItem[K, V]{
		key:   key,
		value: value,
	}
	return true
}

// delete deletes an item and returns the value, or false if not found.
func (n *leafNode[K, V]) delete(key K, compare compare.Compare[K]) (zero V, _ bool) {
	if !n.mut {
		panic("operation on immutable node")
	}

	// Find item by key
	index := n.indexOf(key, compare)

	// Return if not found
	if index >= len(n.items) {
		return zero, false
	}
	if compare(n.items[index].key, key) != 0 {
		return zero, false
	}

	// Get value
	item := n.items[index]
	value := item.value

	// Shift greater items left
	copy(n.items[index:], n.items[index+1:])
	n.items[len(n.items)-1] = leafItem[K, V]{}

	// Truncate items
	n.items = n.items[:len(n.items)-1]
	return value, true
}

// contains/indexOf

// contains returns true if the key exists.
func (n *leafNode[K, V]) contains(key K, compare compare.Compare[K]) bool {
	index := n.indexOf(key, compare)
	if index >= len(n.items) {
		return false
	}

	cmp := compare(n.items[index].key, key)
	return cmp == 0
}

// indexOf returns an index of an item with key >= key, or -1 if not found.
func (n *leafNode[K, V]) indexOf(key K, compare compare.Compare[K]) int {
	return sort.Search(len(n.items), func(i int) bool {
		key0 := n.items[i].key
		cmp := compare(key0, key)
		return cmp >= 0
	})
}

// clone

// clone returns a mutable copy, retains the children.
func (n *leafNode[K, V]) clone() node[K, V] {
	return cloneLeafNode(n)
}

// freeze makes the node immutable.
func (n *leafNode[K, V]) freeze() {
	n.mut = false
}

// split

// split splits the node, and returns the new node, or false if no split required.
func (n *leafNode[K, V]) split() (node[K, V], bool) {
	if !n.mut {
		panic("operation on immutable node")
	}

	if len(n.items) < maxItems {
		return nil, false
	}

	// Calc middle index
	middle := len(n.items) / 2

	// Move items to next node
	next := nextLeafNode(n.items[middle:len(n.items)])

	// Clear and truncate items
	for i := middle; i < len(n.items); i++ {
		n.items[i] = leafItem[K, V]{}
	}
	n.items = n.items[:middle]
	return next, true
}

// refs

// retain increments the reference count.
func (n *leafNode[K, V]) retain() {
	v := atomic.AddInt64(&n.refs, 1)
	if v == 1 {
		panic("retained already released node")
	}
}

// release decrements the reference count and frees the node if the count is zero.
func (n *leafNode[K, V]) release() {
	v := atomic.AddInt64(&n.refs, -1)
	switch {
	case v < 0:
		panic("released already released node")
	case v > 0:
		return
	}

	// Clear items
	n.items = slices2.Truncate(n.items)

	// Release node
	releaseLeaf[K, V](n)
}

// refcount returns the reference count.
func (n *leafNode[K, V]) refcount() int64 {
	return n.refs
}

// pool

var leafNodePools = pools.NewPools()

func acquireLeaf[K, V any]() *leafNode[K, V] {
	v, ok := pools.Acquire[*leafNode[K, V]](leafNodePools)
	if ok {
		return v
	}
	return &leafNode[K, V]{}
}

func releaseLeaf[K, V any](n *leafNode[K, V]) {
	n.reset()
	pools.Release(leafNodePools, n)
}

func (n *leafNode[K, V]) reset() {
	n.items = slices2.Truncate(n.items)

	*n = leafNode[K, V]{}
}
