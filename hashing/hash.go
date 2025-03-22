// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package hashing

import (
	"fmt"
	"math"
	"unsafe"

	"github.com/basecomplextech/baselibrary/bin"
)

// HashBool returns a hash of a bool value.
func HashBool(b bool) uint32 {
	if b {
		return 1
	}
	return 0
}

// HashBytes returns a hash of a byte slice value.
func HashBytes(b []byte) uint32 {
	h := uint32(0)
	for _, c := range b {
		h = h*31 + uint32(c)
	}
	return h
}

// HashString returns a hash of a string value.
func HashString(s string) uint32 {
	h := uint32(0)
	for _, c := range s {
		h = h*31 + uint32(c)
	}
	return h
}

// HashPointer returns a hash of a pointer value.
func HashPointer(v any) uint32 {
	if v == nil {
		return 0
	}

	p := *(*unsafe.Pointer)(unsafe.Pointer(&v))
	h := uintptr(p)
	return uint32(h ^ (h >> 32)) // xor of two halves
}

// Hash returns a hash of a key, panics if the key type is not supported.
// Add more types as needed.
func Hash[K any](key K) uint32 {
	switch v := any(key).(type) {
	case bool:
		if v {
			return 1
		}
		return 0

	case int:
		return uint32(v)

	case int8:
		return uint32(v)
	case int16:
		return uint32(v)
	case int32:
		return uint32(v)
	case int64:
		return uint32(v ^ (v >> 32)) // xor of two halves

	case uint8:
		return uint32(v)
	case uint16:
		return uint32(v)
	case uint32:
		return uint32(v)
	case uint64:
		return uint32(v ^ (v >> 32)) // xor of two halves

	case float32:
		return math.Float32bits(v)
	case float64:
		v1 := math.Float64bits(v)
		return uint32(v1 ^ (v1 >> 32)) // xor of two halves

	case string:
		return HashString(v)
	case []byte:
		return HashBytes(v)

	case bin.Bin64:
		return v.Hash32()
	case bin.Bin128:
		return v.Hash32()
	case bin.Bin256:
		return v.Hash32()
	}

	if h, ok := ((any)(key)).(Hasher); ok {
		return h.Hash32()
	}

	panic(fmt.Sprintf("unsupported type %T", key))
}
