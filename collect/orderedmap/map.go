package orderedmap

import (
	"github.com/basecomplextech/baselibrary/collect/maps"
	"github.com/basecomplextech/baselibrary/collect/slices"
)

// Map is an ordered map which maintains the order of insertion, even on updates.
type Map[K comparable, V any] struct {
	list []Item[K, V] // array list
	map_ map[K]int    // map to index in array list
}

// Item is a key-value pair.
type Item[K comparable, V any] struct {
	Key   K
	Value V
}

// New returns a new ordered map.
func New[K comparable, V any](items ...Item[K, V]) *Map[K, V] {
	m := &Map[K, V]{
		list: make([]Item[K, V], 0, len(items)),
		map_: make(map[K]int, len(items)),
	}

	for _, item := range items {
		m.Put(item.Key, item.Value)
	}
	return m
}

// Index returns the index of the given key, or -1.
func (m *Map[K, V]) Index(key K) int {
	i, ok := m.map_[key]
	if !ok {
		return -1
	}
	return i
}

// Len returns the number of items in the map.
func (m *Map[K, V]) Len() int {
	return len(m.map_)
}

// Get returns the value for the given key, or false.
func (m *Map[K, V]) Get(key K) (value V, ok bool) {
	i, ok := m.map_[key]
	if !ok {
		return value, false
	}

	item := m.list[i]
	return item.Value, true
}

// Put adds or updates the given key with the given value.
func (m *Map[K, V]) Put(key K, value V) {
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
func (m *Map[K, V]) Delete(key K) {
	i, ok := m.map_[key]
	if !ok {
		return
	}

	// Delete item, left shift others
	delete(m.map_, key)
	copy(m.list[i:], m.list[i+1:])
	m.list = m.list[:len(m.list)-1]
}

// Clear removes all items from the map.
func (m *Map[K, V]) Clear() {
	m.list = slices.Clear(m.list)
	maps.Clear(m.map_)
}

// Clone returns a clone of the map.
func (m *Map[K, V]) Clone() *Map[K, V] {
	m1 := New[K, V]()
	for _, item := range m.list {
		m1.Put(item.Key, item.Value)
	}
	return m1
}

// Iterate iterates over the map, returns false if the iteration stopped early.
func (m *Map[K, V]) Iterate(yield func(key K, value V) bool) bool {
	for _, item := range m.list {
		if !yield(item.Key, item.Value) {
			return false
		}
	}
	return true
}

// Keys returns a slice of keys in the order they were inserted.
func (m *Map[K, V]) Keys() []K {
	keys := make([]K, len(m.list))
	for i, item := range m.list {
		keys[i] = item.Key
	}
	return keys
}

// Values returns a slice of values in the order they were inserted.
func (m *Map[K, V]) Values() []V {
	values := make([]V, len(m.list))
	for i, item := range m.list {
		values[i] = item.Value
	}
	return values
}
