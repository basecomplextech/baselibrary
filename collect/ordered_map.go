// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package collect

import (
	"github.com/basecomplextech/baselibrary/collect/internal/orderedmap"
)

// OrderedMap is a map that maintains the insertion order of items.
type OrderedMap[K comparable, V any] = orderedmap.OrderedMap[K, V]

// NewOrderedMap returns a new ordered map.
func NewOrderedMap[K comparable, V any]() OrderedMap[K, V] {
	return orderedmap.New[K, V]()
}

// NewOrderedMapCap returns a new ordered map with the given capacity.
func NewOrderedMapCap[K comparable, V any](cap int) OrderedMap[K, V] {
	return orderedmap.NewCap[K, V](cap)
}
