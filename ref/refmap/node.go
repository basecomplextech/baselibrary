// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package refmap

import (
	"github.com/basecomplextech/baselibrary/ref"
)

type node[K, V any] interface {
	retain()
	release()
	refcount() int64

	clone() node[K, V]
	freeze()
	mutable() bool

	length() int
	minKey() K
	maxKey() K

	indexOf(key K, cmp CompareFunc[K]) int
	get(key K, cmp CompareFunc[K]) (ref.R[V], bool)
	put(key K, value ref.R[V], cmp CompareFunc[K]) bool
	delete(key K, cmp CompareFunc[K]) bool
	contains(key K, cmp CompareFunc[K]) bool
	split() (node[K, V], bool)
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
