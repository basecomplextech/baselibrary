// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package collect

import "github.com/basecomplextech/baselibrary/collect/internal/orderedmap"

type (
	// Map is an ordered map which maintains the order of insertion, even on updates.
	OrderedMap[K comparable, V any] = orderedmap.OrderedMap[K, V]

	// OrderedMapItem is a key-value pair.
	OrderedMapItem[K comparable, V any] = orderedmap.Item[K, V]
)

// NewOrderedMap returns a new ordered map.
func NewOrderedMap[K comparable, V any](items ...OrderedMapItem[K, V]) OrderedMap[K, V] {
	return orderedmap.New(items...)
}

// NewOrderedMapSize returns a new ordered map with the given size hint.
func NewOrderedMapSize[K comparable, V any](size int) OrderedMap[K, V] {
	return orderedmap.NewSize[K, V](size)
}
