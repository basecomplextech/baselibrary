// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package refmap

import "github.com/basecomplextech/baselibrary/ref"

// items returns items as a slice.
func (t *btree[K, V]) items() testItems[K, V] {
	result := make(testItems[K, V], 0, t.length)

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
				item1 := testItem[K, V]{
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
func (t *btree[K, V]) values() []ref.R[V] {
	result := make([]ref.R[V], 0, t.length)

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
