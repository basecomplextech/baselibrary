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
		return uint32(v ^ (v >> 32)) // xor of two halves
	case int8:
		return uint32(v)
	case int16:
		return uint32(v)
	case int32:
		return uint32(v)
	case int64:
		return uint32(v ^ (v >> 32)) // xor of two halves

	case uint:
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

	if h, ok := ((any)(key)).(Hashable); ok {
		return h.Hash32()
	}

	panic(fmt.Sprintf("unsupported type %T", key))
}

// HashBool returns a hash of a bool value.
func HashBool(b bool) uint32 {
	if b {
		return 1
	}
	return 0
}

// HashInt returns a hash of an int value.
func HashInt(i int) uint32     { return uint32(i ^ (i >> 32)) } // xor of two halves
func HashInt8(i int8) uint32   { return uint32(i) }
func HashInt16(i int16) uint32 { return uint32(i) }
func HashInt32(i int32) uint32 { return uint32(i) }
func HashInt64(i int64) uint32 { return uint32(i ^ (i >> 32)) } // xor of two halves

// HashUint returns a hash of a uint value.
func HashUint(u uint) uint32     { return uint32(u ^ (u >> 32)) } // xor of two halves
func HashUint8(u uint8) uint32   { return uint32(u) }
func HashUint16(u uint16) uint32 { return uint32(u) }
func HashUint32(u uint32) uint32 { return u }
func HashUint64(u uint64) uint32 { return uint32(u ^ (u >> 32)) } // xor of two halves

// HashFloat returns a hash of a float value.
func HashFloat32(f float32) uint32 { return math.Float32bits(f) }
func HashFloat64(f float64) uint32 {
	v := math.Float64bits(f)
	return uint32(v ^ (v >> 32)) // xor of two halves
}

// HashBin returns a hash of a bin value.
func HashBin64(b bin.Bin64) uint32   { return b.Hash32() }
func HashBin128(b bin.Bin128) uint32 { return b.Hash32() }
func HashBin256(b bin.Bin256) uint32 { return b.Hash32() }

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
// The method is unsafe and must be used only with actual pointers.
func HashPointer(v any) uint32 {
	type eface struct {
		typ unsafe.Pointer
		ptr unsafe.Pointer
	}

	if v == nil {
		return 0
	}

	e := (*eface)(unsafe.Pointer(&v))
	if e.typ == nil {
		return 0
	}

	h := uintptr(e.ptr)
	return uint32(h ^ (h >> 32)) // xor of two halves
}
