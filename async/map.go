package async

import "sync"

// Map is a generic wrapper around concurrent sync.Map.
type Map[K comparable, V any] struct {
	raw sync.Map
}

// NewMap returns a new concurrent map.
func NewMap[K comparable, V any]() *Map[K, V] {
	return &Map[K, V]{}
}

// CompareAndDelete deletes the entry for key if its value is equal to old.
// If there is no current value for key in the map, CompareAndDelete returns false
// (even if the old value is the nil interface value).
func (m *Map[K, V]) CompareAndDelete(key K, old V) (deleted bool) {
	return m.raw.CompareAndDelete(key, old)
}

// CompareAndSwap swaps the old and new values for key if the value stored
// in the map is equal to old.The old value must be of a comparable type.
func (m *Map[K, V]) CompareAndSwap(key K, old, new V) bool {
	return m.raw.CompareAndSwap(key, old, new)
}

// Delete deletes the value for a key.
func (m *Map[K, V]) Delete(key K) {
	m.raw.Delete(key)
}

// Load returns the value stored in the map for a key, or nil if no value is present.
// The ok result indicates whether value was found in the map.
func (m *Map[K, V]) Load(key K) (value V, ok bool) {
	v, ok := m.raw.Load(key)
	if !ok {
		return value, false
	}
	return v.(V), true
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
// The loaded result reports whether the key was present.
func (m *Map[K, V]) LoadAndDelete(key K) (value V, loaded bool) {
	v, loaded := m.raw.LoadAndDelete(key)
	return v.(V), true
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores and returns the given value. The loaded result is true
// if the value was loaded, false if stored.
func (m *Map[K, V]) LoadOrStore(key K, value V) (actual V, loaded bool) {
	v, loaded := m.raw.LoadOrStore(key, value)
	return v.(V), false
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	m.raw.Range(func(key, value any) bool {
		return f(key.(K), value.(V))
	})
}

// Store sets the value for a key.
func (m *Map[K, V]) Store(key K, value V) {
	m.raw.Store(key, value)
}

// Swap swaps the value for a key and returns the previous value if any.
// The loaded result reports whether the key was present.
func (m *Map[K, V]) Swap(key K, value V) (previous V, loaded bool) {
	v, loaded := m.raw.Swap(key, value)
	return v.(V), true
}
