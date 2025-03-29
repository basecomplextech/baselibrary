// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package orderedmap

import (
	"slices"

	"github.com/basecomplextech/baselibrary/collect/slices2"
)

// OrderedMap is a map that maintains the insertion order of items.
type OrderedMap[K comparable, V any] interface {
	// Len returns the number of items in the map.
	Len() int

	// Read

	// Get returns a value by a key, or false.
	Get(key K) (V, bool)

	// GetAt returns an item at a given index.
	GetAt(index int) (K, V)

	// Contains returns true if a key exists.
	Contains(key K) bool

	// IndexOf returns the index of a key, or -1.
	// The method takes O(n) time in the worst case.
	IndexOf(key K) int

	// Keys returns a slice of keys in the order they were inserted.
	Keys() []K

	// Values returns a slice of values in the order they were inserted.
	Values() []V

	// Write

	// Clear deletes all items.
	Clear()

	// Clone returns a copy of the map.
	Clone() OrderedMap[K, V]

	// Delete deletes an item, and returns its value.
	// The method takes O(n) time in the worst case.
	Delete(key K) (V, bool)

	// Set sets an item.
	Set(key K, value V)
}

// New returns a new ordered map.
func New[K comparable, V any]() OrderedMap[K, V] {
	return newOrderedMap[K, V]()
}

// NewCap returns a new ordered map with a given capacity.
func NewCap[K comparable, V any](cap int) OrderedMap[K, V] {
	return newOrderedMapCap[K, V](cap)
}

// internal

var _ OrderedMap[int, any] = (*orderedMap[int, any])(nil)

type orderedMap[K comparable, V any] struct {
	items  []item[K, V]
	values map[K]V
}

type item[K comparable, V any] struct {
	key   K
	value V
}

func newOrderedMap[K comparable, V any](items ...item[K, V]) *orderedMap[K, V] {
	m := &orderedMap[K, V]{
		values: make(map[K]V, len(items)),
	}
	m.items = slices.Grow(m.items, len(items))

	for _, item := range items {
		m.Set(item.key, item.value)
	}
	return m
}

func newOrderedMapCap[K comparable, V any](cap int) *orderedMap[K, V] {
	return &orderedMap[K, V]{
		values: make(map[K]V, cap),
		items:  make([]item[K, V], 0, cap),
	}
}

// Len returns the number of items in the map.
func (m *orderedMap[K, V]) Len() int {
	return len(m.items)
}

// Read

// Get returns a value by a key, or false.
func (m *orderedMap[K, V]) Get(key K) (V, bool) {
	v, ok := m.values[key]
	return v, ok
}

// GetAt returns an item at a given index.
func (m *orderedMap[K, V]) GetAt(index int) (K, V) {
	item := m.items[index]
	return item.key, item.value
}

// Contains returns true if a key exists.
func (m *orderedMap[K, V]) Contains(key K) bool {
	_, ok := m.values[key]
	return ok
}

// IndexOf returns the index of a key, or -1.
// The method takes O(n) time in the worst case.
func (m *orderedMap[K, V]) IndexOf(key K) int {
	return m.indexOf(key)
}

// Keys returns a slice of keys in the order they were inserted.
func (m *orderedMap[K, V]) Keys() []K {
	keys := make([]K, len(m.items))
	for i, item := range m.items {
		keys[i] = item.key
	}
	return keys
}

// Values returns a slice of values in the order they were inserted.
func (m *orderedMap[K, V]) Values() []V {
	values := make([]V, len(m.items))
	for i, item := range m.items {
		values[i] = item.value
	}
	return values
}

// Write

// Clear deletes all items.
func (m *orderedMap[K, V]) Clear() {
	m.items = slices2.Truncate(m.items)
	clear(m.values)
}

// Clone returns a copy of the map.
func (m *orderedMap[K, V]) Clone() OrderedMap[K, V] {
	m1 := &orderedMap[K, V]{
		items:  slices.Clone(m.items),
		values: make(map[K]V, len(m.values)),
	}

	for k, v := range m.values {
		m1.values[k] = v
	}
	return m1
}

// Delete deletes an item, and returns its value.
// The method takes O(n) time in the worst case.
func (m *orderedMap[K, V]) Delete(key K) (V, bool) {
	v, ok := m.values[key]
	if !ok {
		return v, false
	}

	index := m.indexOf(key)
	delete(m.values, key)
	m.items = slices.Delete(m.items, index, index+1)
	return v, true
}

// Set sets an item.
func (m *orderedMap[K, V]) Set(key K, value V) {
	// Update existing
	if _, ok := m.values[key]; ok {
		index := m.indexOf(key)
		m.items[index].value = value
		m.values[key] = value
		return
	}

	// Insert new
	item := item[K, V]{key: key, value: value}
	m.items = append(m.items, item)
	m.values[key] = value
}

// private

func (m *orderedMap[K, V]) indexOf(key K) int {
	return slices.IndexFunc(m.items, func(item item[K, V]) bool {
		return item.key == key
	})
}
