package maps2

// Clone returns a map clone, skips nil maps.
func Clone[M ~map[K]V, K comparable, V any](m M) M {
	if m == nil {
		return nil
	}

	m1 := make(map[K]V, len(m))
	for k, v := range m {
		m1[k] = v
	}
	return m1
}

// Copy copies items from src to dst.
func Copy[M ~map[K]V, K comparable, V any](dst, src M) {
	for k, v := range src {
		dst[k] = v
	}
}

// Keys returns a slice of map keys, skips nil maps.
func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	if m == nil {
		return nil
	}

	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// Values returns a slice of map values, skips nil maps.
func Values[M ~map[K]V, K comparable, V any](m M) []V {
	if m == nil {
		return nil
	}

	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}
