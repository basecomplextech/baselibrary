// Copyright 2023 Ivan Korobkov. All rights reserved.

package refmap

import (
	"github.com/basecomplextech/baselibrary/collect/slices2"
	"github.com/basecomplextech/baselibrary/pools"
	"github.com/basecomplextech/baselibrary/ref"
	"github.com/basecomplextech/baselibrary/status"
)

// Iterator sequentially iterates over sorted map items.
//
// Usage:
//
//	it := refmap.Iterator()
//	defer it.Free()
//
//	it.SeekToStart()
//
//	for it.Next() {
//		key := it.Key()
//		value := it.Value()
//	}
type Iterator[K any, V any] interface {
	// OK returns true when the iterator points to a valid item, or false on end.
	OK() bool

	// Key returns the current key or zero, the key is valid until the next iteration.
	Key() K

	// Value returns the current value or zero, the value is valid until the next iteration.
	Value() ref.R[V]

	// Iterating

	// Next moves to the next item.
	Next() bool

	// Previous moves to the previous item.
	Previous() bool

	// Seeking

	// SeekToStart positions the iterator at the start.
	SeekToStart() bool

	// SeekToEnd positions the iterator at the end.
	SeekToEnd() bool

	// SeekBefore positions the iterator before an item with key >= key.
	SeekBefore(key K) bool

	// Internal

	// Free frees the iterator, implements the ref.Free interface.
	Free()
}

// internal

type position int

const (
	positionUndefined position = iota
	positionItem
	positionBefore
	positionStart
	positionEnd
)

var _ Iterator[int, int] = (*iterator[int, int])(nil)

// iterator iterates over a btree, does not retain the values.
type iterator[K, V any] struct {
	*iterState[K, V]
}

type iterState[K, V any] struct {
	tree *btree[K, V]

	st    status.Status    // current iterator status
	mod   int              // track concurrent modifications
	pos   position         // current iterator position
	stack []iterElem[K, V] // the last element is always a leaf node
}

// iterElem combines a node and an iteration index of the node elements.
type iterElem[K, V any] struct {
	node[K, V]
	index int // index maybe == len(node.items) as in sort.Search
}

func newIterator[K, V any](tree *btree[K, V]) *iterator[K, V] {
	it := &iterator[K, V]{acquireIterState[K, V]()}
	it.tree = tree

	it.st = status.None
	it.mod = tree.mod
	return it
}

// reset

func (s *iterState[K, V]) reset() {
	stack := slices2.Truncate(s.stack)

	*s = iterState[K, V]{}
	s.stack = stack
}

// State

// OK returns true when the iterator points to a valid item, or false on end.
func (it *iterator[K, V]) OK() bool {
	return it.st.OK()
}

// Key returns the current key or zero, the key is valid until the next iteration.
func (it *iterator[K, V]) Key() (key K) {
	if !it.st.OK() {
		return
	}

	elem := it.stack[len(it.stack)-1]
	node := elem.node.(*leafNode[K, V])
	index := elem.index

	if index >= len(node.items) {
		return
	}
	return node.items[index].key
}

// Value returns the current value or zero, the value is valid until the next iteration.
func (it *iterator[K, V]) Value() (value ref.R[V]) {
	if !it.st.OK() {
		return
	}

	elem := it.stack[len(it.stack)-1]
	node := elem.node.(*leafNode[K, V])
	index := elem.index

	if index >= len(node.items) {
		return
	}
	return node.items[index].value
}

// Iterating

// Next moves to the next item.
func (it *iterator[K, V]) Next() bool {
	switch it.st.Code {
	case status.CodeOK,
		status.CodeEnd,
		status.CodeNone:
	default:
		return false
	}

	if it.mod != it.tree.mod {
		it.st = status.Errorf("refmap concurrent modification error")
		panic("refmap concurrent modification error")
	}

	switch it.pos {
	case positionBefore:
		// Stack already points to the next item
		// Just change the state
		it.pos = positionItem
		it.st = status.OK
		return true

	case positionStart:
		// Recursively push first nodes
		// And check stack is not empty
		it.stack = it.stack[:0]
		it.pushStart(it.tree.root)

		if len(it.stack) == 0 {
			it.pos = positionEnd
			it.st = status.End
			return false
		}

		it.pos = positionItem
		it.st = status.OK
		return true

	case positionEnd:
		// Cannot proceed
		it.st = status.End
		return false

	case positionItem:
		// Move to the next item

		for len(it.stack) > 0 {
			// Peek top node
			elem := &it.stack[len(it.stack)-1]
			node := elem.node
			last := node.length() - 1

			// Pop it and continue on end
			if elem.index == last {
				it.stack = it.stack[:len(it.stack)-1]
				continue
			}

			// Move to next element,
			// And break if leaf node
			elem.index++
			n, ok := node.(*branchNode[K, V])
			if !ok {
				break
			}

			// Recursively push children onto the stack
			child := n.items[elem.index].node
			it.pushStart(child)
			break
		}

		if len(it.stack) == 0 {
			it.pos = positionEnd
			it.st = status.End
			return false
		}

		it.pos = positionItem
		it.st = status.OK
		return true
	}

	it.st = status.End
	return false
}

// Previous moves to the previous item.
func (it *iterator[K, V]) Previous() bool {
	switch it.st.Code {
	case status.CodeOK,
		status.CodeEnd,
		status.CodeNone:
	default:
		return false
	}

	if it.mod != it.tree.mod {
		it.st = status.Errorf("refmap concurrent modification error")
		panic("refmap concurrent modification error")
	}

	switch it.pos {
	case positionStart:
		// Cannot proceed
		it.st = status.End
		return false

	case positionEnd:
		// Recursively push last nodes
		// And check stack is not empty
		it.stack = it.stack[:0]
		it.pushEnd(it.tree.root)

		if len(it.stack) == 0 {
			it.pos = positionStart
			it.st = status.End
			return false
		}

		it.pos = positionItem
		it.st = status.OK
		return true

	case positionBefore, positionItem:
		// Move to the previous item

		for len(it.stack) > 0 {
			// Peek top node
			elem := &it.stack[len(it.stack)-1]
			node := elem.node

			// Pop it and continue on start
			if elem.index == 0 {
				it.stack = it.stack[:len(it.stack)-1]
				continue
			}

			// Move to previous element,
			// And break if leaf node
			elem.index--
			n, ok := node.(*branchNode[K, V])
			if !ok {
				break
			}

			// Recursively push children onto the stack
			child := n.items[elem.index].node
			it.pushEnd(child)
			break
		}

		if len(it.stack) == 0 {
			it.pos = positionStart
			it.st = status.End
			return false
		}

		it.pos = positionItem
		it.st = status.OK
		return true
	}

	it.st = status.End
	return false
}

// Seeking

// SeekToStart positions the iterator at the start.
func (it *iterator[K, V]) SeekToStart() bool {
	switch it.st.Code {
	case status.CodeOK,
		status.CodeEnd,
		status.CodeNone:
	default:
		return false
	}

	it.st = status.None
	it.mod = it.tree.mod
	it.pos = positionStart
	it.stack = it.stack[:0]
	return true
}

// SeekToEnd positions the iterator at the end.
func (it *iterator[K, V]) SeekToEnd() bool {
	switch it.st.Code {
	case status.CodeOK,
		status.CodeEnd,
		status.CodeNone:
	default:
		return false
	}

	it.st = status.None
	it.mod = it.tree.mod
	it.pos = positionEnd
	it.stack = it.stack[:0]
	return true
}

// SeekBefore positions the iterator before an item with key >= key, and returns ok/end/error.
func (it *iterator[K, V]) SeekBefore(key K) bool {
	switch it.st.Code {
	case status.CodeOK,
		status.CodeEnd,
		status.CodeNone:
	default:
		return false
	}

	it.st = status.None
	it.stack = it.stack[:0]

	// Recursively push nodes onto the stack
	// And position them at elements >= key
	node := it.tree.root
	for node != nil {
		index := node.indexOf(key, it.tree.compare)
		elem := iterElem[K, V]{
			node:  node,
			index: index,
		}

		// Push element onto the stack
		// If this is a branch node
		branch, ok := node.(*branchNode[K, V])
		if ok {
			node = branch.child(index)
			it.stack = append(it.stack, elem)
			continue
		}

		// Check if a leaf node contains a key >= given key,
		// Otherwise clear the stack to end iteration
		leaf := node.(*leafNode[K, V])
		if index >= len(leaf.items) {
			it.stack = it.stack[:0]
			break
		}

		it.stack = append(it.stack, elem)
		break
	}

	if len(it.stack) == 0 {
		it.st = status.End
		it.mod = it.tree.mod
		it.pos = positionEnd
		return false
	}

	it.st = status.None
	it.mod = it.tree.mod
	it.pos = positionBefore
	return true

}

// Internal

// Free frees the iterator.
func (it *iterator[K, V]) Free() {
	state := it.iterState
	it.iterState = nil
	releaseIterState(state)
}

// private

// pushStart recursively pushes nodes onto the stack
// and positions them at start elements.
func (it *iterator[K, V]) pushStart(node node[K, V]) {
	for node != nil {
		if node.length() == 0 {
			break
		}

		index := 0
		elem := iterElem[K, V]{
			node:  node,
			index: index,
		}
		it.stack = append(it.stack, elem)

		n, ok := node.(*branchNode[K, V])
		if !ok {
			break
		}

		node = n.child(index)
	}
}

// pushEnd recursively pushes nodes onto the stack
// and positions them at end elements.
func (it *iterator[K, V]) pushEnd(node node[K, V]) {
	for node != nil {
		if node.length() == 0 {
			break
		}

		index := node.length() - 1
		elem := iterElem[K, V]{
			node:  node,
			index: index,
		}
		it.stack = append(it.stack, elem)

		n, ok := node.(*branchNode[K, V])
		if !ok {
			break
		}

		node = n.child(index)
	}
}

// iterator state pool

var iterStatePools = pools.NewPools()

func acquireIterState[K, V any]() *iterState[K, V] {
	s, ok := pools.Acquire[*iterState[K, V]](iterStatePools)
	if ok {
		return s
	}
	return &iterState[K, V]{}
}

func releaseIterState[K, V any](s *iterState[K, V]) {
	s.reset()
	pools.Release(iterStatePools, s)
}
