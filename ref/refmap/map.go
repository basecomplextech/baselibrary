package refmap

import (
	"github.com/basecomplextech/baselibrary/compare"
	"github.com/basecomplextech/baselibrary/ref"
)

// Map is an immutable sorted map which stores countable references, implemented as a btree.
//
// The map retains/releases references internally, but does not retain them when iterating
// or returning.
type Map[K any, V ref.Ref] interface {
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

	// Contains returns true if the map contains a key.
	Contains(key K) bool

	// Iterator returns an iterator positioned at the start of the map.
	Iterator() Iterator[K, V]

	// Write

	// Put adds an item to the map.
	Put(key K, value V)

	// Delete deletes an item by a key.
	Delete(key K)

	// Internal

	// Free frees the map, implements the ref.Freer interface.
	Free()
}

// Item is a map item.
type Item[K any, V ref.Ref] struct {
	Key   K
	Value V
}

// New returns an empty map.
func New[K any, V ref.Ref](mutable bool, compare compare.Func[K]) Map[K, V] {
	m := newBtree[K, V](compare)
	if !mutable {
		m.Freeze()
	}
	return m
}

// New returns an empty map wrapped in a ref.
func NewRef[K any, V ref.Ref](mutable bool, compare compare.Func[K]) *ref.R[Map[K, V]] {
	m := New[K, V](mutable, compare)
	return ref.New(m)
}

// internal

const maxItems = 16

var _ Map[any, ref.Ref] = (*btree[any, ref.Ref])(nil)

type btree[K any, V ref.Ref] struct {
	compare compare.Func[K]

	root    node[K, V]
	height  int
	length  int64
	mutable bool
}

func newBtree[K any, V ref.Ref](compare compare.Func[K]) *btree[K, V] {
	return &btree[K, V]{
		compare: compare,
		mutable: true,

		root:   newLeafNode[K, V](),
		height: 1,
	}
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
	t1 := &btree[K, V]{
		compare: t.compare,

		root:    root1,
		height:  t.height,
		length:  t.length,
		mutable: true,
	}
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

// Get returns a value by a key, does not retain the value.
func (t *btree[K, V]) Get(key K) (v V, ok bool) {
	return t.root.get(key, t.compare)
}

// Contains returns true if the map contains a key.
func (t *btree[K, V]) Contains(key K) bool {
	return t.root.contains(key, t.compare)
}

// Iterator returns an iterator, the iterator does not retain the values.
func (t *btree[K, V]) Iterator() Iterator[K, V] {
	it := newIterator(t)
	it.SeekToStart()
	return it
}

// Write

// Put adds an item to the map.
func (t *btree[K, V]) Put(key K, value V) {
	if !t.mutable {
		panic("operation on immutable refmap")
	}

	// Split if full
	t.maybeSplitRoot()

	// Insert item
	node := t.mutateRoot()
	mod := node.put(key, value, t.compare)
	if !mod {
		return
	}

	t.length++
}

// Delete deletes an item by a key.
func (t *btree[K, V]) Delete(key K) {
	if !t.mutable {
		panic("operation on immutable refmap")
	}

	// Mutate root
	node := t.mutateRoot()

	// Delete item
	mod := node.delete(key, t.compare)
	if !mod {
		return
	}

	// Decrement length
	t.length--
}

// Internal

// Free frees the map, implements the ref.Freer interface.
func (t *btree[K, V]) Free() {
	t.root.release()

	t.root = nil
	t.height = 0
	t.length = 0
}

// root

func (t *btree[K, V]) mutateRoot() node[K, V] {
	if !t.mutable {
		panic("operation on immutable refmap")
	}

	if t.root.mutable() {
		return t.root
	}
	prev := t.root

	// Clone and replace root
	next := t.root.clone()
	t.root = next

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
