// Copyright 2024 Ivan Korobkov. All rights reserved.

package hashing

import (
	"fmt"
	"math"

	"github.com/basecomplextech/baselibrary/bin"
)

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
		h := uint32(0)
		for _, c := range v {
			h = h*31 + uint32(c)
		}
		return h

	case []byte:
		h := uint32(0)
		for _, c := range v {
			h = h*31 + uint32(c)
		}
		return h

	case bin.Bin64:
		return v.Hash32()
	case bin.Bin128:
		return v.Hash32()
	case bin.Bin256:
		return v.Hash32()
	}

	panic(fmt.Sprintf("unsupported type %T", key))
}
