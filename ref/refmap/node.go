// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package refmap

import (
	"github.com/basecomplextech/baselibrary/ref"
)

type node[K, V any] interface {
	// length returns the number of items in the node.
	length() int

	// minKey returns the minimum key in the node.
	minKey() K

	// maxKey returns the maximum key in the node.
	maxKey() K

	// mutable returns true if the node is mutable.
	mutable() bool

	// get/insert/delete

	// get returns for an item by key, or false if not found.
	get(key K, cmp CompareFunc[K]) (ref.R[V], bool)

	// insert inserts or updates an item, returns true if inserted.
	insert(key K, value ref.R[V], cmp CompareFunc[K]) bool

	// delete deletes an item by key, returns true if deleted.
	delete(key K, cmp CompareFunc[K]) bool

	// contains/indexOf

	// contains returns true if the key exists.
	contains(key K, cmp CompareFunc[K]) bool

	// indexOf returns an index of an item with key >= key, or -1 if not found.
	indexOf(key K, cmp CompareFunc[K]) int

	// clone

	// clone returns a mutable copy, retains the children.
	clone() node[K, V]

	// freeze makes the node immutable.
	freeze()

	// split

	// split splits the node, and returns the new node, or false if no split required.
	split() (node[K, V], bool)

	// refs

	// retain increments the reference count.
	retain()

	// release decrements the reference count and frees the node if the count is zero.
	release()

	// refcount returns the reference count.
	refcount() int64
}

func walk[K, V any](node node[K, V], fn func(node[K, V])) {
	fn(node)

	n, ok := node.(*branchNode[K, V])
	if !ok {
		return
	}

	for _, item := range n.items {
		walk(item.node, fn)
	}
}
