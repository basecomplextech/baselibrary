// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package immap

import (
	"github.com/basecomplextech/baselibrary/collect/slices2"
	"github.com/basecomplextech/baselibrary/iterator/mapiter"
	"github.com/basecomplextech/baselibrary/pools"
)

// Iterator iterates over an immutable map, extends mapiter.Iter with reverse iteration and seeking.
type Iterator[K any, V any] interface {
	mapiter.Iter[K, V]

	// Next returns the next key-value pair, or false on the end.
	Next() (K, V, bool)

	// Previous returns the previous key-value pair, or false on the end.
	Previous() (K, V, bool)

	// Seeking

	// SeekToStart positions the iterator at the start.
	SeekToStart()

	// SeekToEnd positions the iterator at the end.
	SeekToEnd()

	// SeekBefore positions the iterator before an item with key >= key, or false on the end.
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

var _ Iterator[int, int] = (*iter[int, int])(nil)

// iter iterates over a btree, does not retain the values.
type iter[K, V any] struct {
	*iterState[K, V]
}

type iterState[K, V any] struct {
	tree *btree[K, V]

	mod   int              // track concurrent modifications
	pos   position         // current iterator position
	stack []iterElem[K, V] // the last element is always a leaf node
}

// iterElem combines a node and an iteration index of the node elements.
type iterElem[K, V any] struct {
	node[K, V]
	index int // index maybe == len(node.items) as in sort.Search
}

func newIterator[K, V any](tree *btree[K, V]) *iter[K, V] {
	it := &iter[K, V]{acquireIterState[K, V]()}
	it.tree = tree
	it.mod = tree.mod
	it.pos = positionUndefined
	return it
}

// Next returns the next key-value pair, or false on the end.
func (it *iter[K, V]) Next() (key K, value V, _ bool) {
	if it.mod != it.tree.mod {
		panic("immap concurrent modification error")
	}

	switch it.pos {
	case positionBefore:
		// Stack already points to the next item
		// Just change the state
		it.pos = positionItem

	case positionStart:
		// Recursively push first nodes
		// And check stack is not empty
		it.stack = it.stack[:0]
		it.pushStart(it.tree.root)

		if len(it.stack) == 0 {
			it.pos = positionEnd
			return key, value, false
		}

		it.pos = positionItem

	case positionEnd:
		// Cannot proceed
		return key, value, false

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
			return key, value, false
		}
		it.pos = positionItem
	}

	if len(it.stack) == 0 {
		it.pos = positionEnd
		return key, value, false
	}

	// Return current item
	elem := it.stack[len(it.stack)-1]
	node := elem.node.(*leafNode[K, V])
	if elem.index >= len(node.items) {
		return key, value, false
	}

	item := &node.items[elem.index]
	return item.key, item.value, true
}

// Previous returns the previous key-value pair, or false on the end.
func (it *iter[K, V]) Previous() (key K, value V, _ bool) {
	if it.mod != it.tree.mod {
		panic("immap concurrent modification error")
	}

	switch it.pos {
	case positionStart:
		// Cannot proceed
		return key, value, false

	case positionEnd:
		// Recursively push last nodes
		// And check stack is not empty
		it.stack = it.stack[:0]
		it.pushEnd(it.tree.root)

		if len(it.stack) == 0 {
			it.pos = positionStart
			return key, value, false
		}

		it.pos = positionItem

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
			return key, value, false
		}
		it.pos = positionItem
	}

	if len(it.stack) == 0 {
		it.pos = positionStart
		return key, value, false
	}

	// Return current item
	elem := it.stack[len(it.stack)-1]
	node := elem.node.(*leafNode[K, V])
	if elem.index >= len(node.items) {
		return key, value, false
	}

	item := &node.items[elem.index]
	return item.key, item.value, true
}

// Seeking

// SeekToStart positions the iterator at the start.
func (it *iter[K, V]) SeekToStart() {
	it.mod = it.tree.mod
	it.pos = positionStart
	it.stack = it.stack[:0]
}

// SeekToEnd positions the iterator at the end.
func (it *iter[K, V]) SeekToEnd() {
	it.mod = it.tree.mod
	it.pos = positionEnd
	it.stack = it.stack[:0]
}

// SeekBefore positions the iterator before an item with key >= key, or false on the end.
func (it *iter[K, V]) SeekBefore(key K) bool {
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
		it.mod = it.tree.mod
		it.pos = positionEnd
		return false
	}

	it.mod = it.tree.mod
	it.pos = positionBefore
	return true
}

// Internal

// Free frees the iterator.
func (it *iter[K, V]) Free() {
	state := it.iterState
	it.iterState = nil
	releaseIterState(state)
}

// private

// pushStart recursively pushes nodes onto the stack
// and positions them at start elements.
func (it *iter[K, V]) pushStart(node node[K, V]) {
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
func (it *iter[K, V]) pushEnd(node node[K, V]) {
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

func (s *iterState[K, V]) reset() {
	stack := slices2.Truncate(s.stack)

	*s = iterState[K, V]{}
	s.stack = stack
}
