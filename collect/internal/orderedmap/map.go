// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package orderedmap

import "github.com/basecomplextech/baselibrary/collect/slices2"

// Map is an ordered map which maintains the order of insertion, even on updates.
type OrderedMap[K comparable, V any] interface {
	// Index returns the index of the given key, or -1.
	Index(key K) int

	// Len returns the number of items in the map.
	Len() int

	// Contains returns true if the map contains the given key.
	Contains(key K) bool

	// Get returns the value for the given key, or false.
	Get(key K) (value V, ok bool)

	// Put adds or updates the given key with the given value.
	Put(key K, value V)

	// Delete removes the given key, the method takes O(n) time to update the linked list.
	Delete(key K)

	// Clear removes all items from the map.
	Clear()

	// Clone returns a clone of the map.
	Clone() OrderedMap[K, V]

	// Iterate iterates over the map, returns false if the iteration stopped early.
	Iterate(yield func(key K, value V) bool) bool

	// Item returns an item at the given index.
	Item(index int) (K, V)

	// Key returns a key at the given index.
	Key(index int) K

	// Value returns a value at the given index.
	Value(index int) V

	// Items returns a slice of items in the order they were inserted.
	Items() []Item[K, V]

	// Keys returns a slice of keys in the order they were inserted.
	Keys() []K

	// Values returns a slice of values in the order they were inserted.
	Values() []V
}

// Item is a key-value pair.
type Item[K comparable, V any] struct {
	Key   K
	Value V
}

// New returns a new ordered map.
func New[K comparable, V any](items ...Item[K, V]) OrderedMap[K, V] {
	return newMap[K, V](items...)
}

// NewSize returns a new ordered map with the given size hint.
func NewSize[K comparable, V any](size int) OrderedMap[K, V] {
	return newMapSize[K, V](size)
}

// internal

var _ OrderedMap[int, int] = (*orderedMap[int, int])(nil)

type orderedMap[K comparable, V any] struct {
	list []Item[K, V] // array list
	map_ map[K]int    // map to index in array list
}

// New returns a new ordered map.
func newMap[K comparable, V any](items ...Item[K, V]) *orderedMap[K, V] {
	m := &orderedMap[K, V]{
		list: make([]Item[K, V], 0, len(items)),
		map_: make(map[K]int, len(items)),
	}

	for _, item := range items {
		m.Put(item.Key, item.Value)
	}
	return m
}

// NewSize returns a new ordered map with the given size hint.
func newMapSize[K comparable, V any](size int) *orderedMap[K, V] {
	return &orderedMap[K, V]{
		list: make([]Item[K, V], 0, size),
		map_: make(map[K]int, size),
	}
}

// Index returns the index of the given key, or -1.
func (m *orderedMap[K, V]) Index(key K) int {
	i, ok := m.map_[key]
	if !ok {
		return -1
	}
	return i
}

// Len returns the number of items in the map.
func (m *orderedMap[K, V]) Len() int {
	return len(m.map_)
}

// Contains returns true if the map contains the given key.
func (m *orderedMap[K, V]) Contains(key K) bool {
	_, ok := m.map_[key]
	return ok
}

// Get returns the value for the given key, or false.
func (m *orderedMap[K, V]) Get(key K) (value V, ok bool) {
	i, ok := m.map_[key]
	if !ok {
		return value, false
	}

	item := m.list[i]
	return item.Value, true
}

// Put adds or updates the given key with the given value.
func (m *orderedMap[K, V]) Put(key K, value V) {
	i, ok := m.map_[key]
	if ok {
		m.list[i].Value = value
		return
	}

	i = len(m.list)
	item := Item[K, V]{
		Key:   key,
		Value: value,
	}
	m.list = append(m.list, item)
	m.map_[key] = i
}

// Delete removes the given key, the method takes O(n) time to update the linked list.
func (m *orderedMap[K, V]) Delete(key K) {
	i, ok := m.map_[key]
	if !ok {
		return
	}

	// TODO: Fix this, indexes change

	// Delete item, left shift others
	delete(m.map_, key)
	copy(m.list[i:], m.list[i+1:])
	m.list = m.list[:len(m.list)-1]
}

// Clear removes all items from the map.
func (m *orderedMap[K, V]) Clear() {
	m.list = slices2.Truncate(m.list)
	clear(m.map_)
}

// Clone returns a clone of the map.
func (m *orderedMap[K, V]) Clone() OrderedMap[K, V] {
	m1 := newMapSize[K, V](len(m.list))
	for _, item := range m.list {
		m1.Put(item.Key, item.Value)
	}
	return m1
}

// Iterate iterates over the map, returns false if the iteration stopped early.
func (m *orderedMap[K, V]) Iterate(yield func(key K, value V) bool) bool {
	for _, item := range m.list {
		if !yield(item.Key, item.Value) {
			return false
		}
	}
	return true
}

// Item returns an item at the given index.
func (m *orderedMap[K, V]) Item(index int) (K, V) {
	item := m.list[index]
	return item.Key, item.Value
}

// Key returns a key at the given index.
func (m *orderedMap[K, V]) Key(index int) K {
	return m.list[index].Key
}

// Value returns a value at the given index.
func (m *orderedMap[K, V]) Value(index int) V {
	return m.list[index].Value
}

// Items returns a slice of items in the order they were inserted.
func (m *orderedMap[K, V]) Items() []Item[K, V] {
	items := make([]Item[K, V], len(m.list))
	copy(items, m.list)
	return items
}

// Keys returns a slice of keys in the order they were inserted.
func (m *orderedMap[K, V]) Keys() []K {
	keys := make([]K, len(m.list))
	for i, item := range m.list {
		keys[i] = item.Key
	}
	return keys
}

// Values returns a slice of values in the order they were inserted.
func (m *orderedMap[K, V]) Values() []V {
	values := make([]V, len(m.list))
	for i, item := range m.list {
		values[i] = item.Value
	}
	return values
}
