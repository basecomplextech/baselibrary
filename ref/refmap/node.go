package refmap

import (
	"github.com/basecomplextech/baselibrary/ref"
)

type node[K any, V ref.Ref] interface {
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
	get(key K, cmp CompareFunc[K]) (V, bool)
	put(key K, value V, cmp CompareFunc[K]) bool
	delete(key K, cmp CompareFunc[K]) bool
	contains(key K, cmp CompareFunc[K]) bool
	split() (node[K, V], bool)
}

func walk[K any, V ref.Ref](node node[K, V], fn func(node[K, V])) {
	fn(node)

	n, ok := node.(*branchNode[K, V])
	if !ok {
		return
	}

	for _, item := range n.items {
		walk(item.node, fn)
	}
}
