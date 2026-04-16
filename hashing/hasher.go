// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package hashing

import (
	"fmt"

	"github.com/basecomplextech/baselibrary/bin"
)

// Hasher computes a 32-bit hash of a value.
type Hasher[T any] interface {
	// Hash32 returns a 32-bit hash of the value.
	Hash32(v T) uint32
}

// NewHasher returns a hasher for the given type, panics if the type is not supported.
func NewHasher[T any]() Hasher[T] {
	var k T

	switch any(k).(type) {
	case bool:
		return (any)(boolHasher{}).(Hasher[T])

	case int:
		return (any)(intHasher{}).(Hasher[T])
	case int8:
		return (any)(int8Hasher{}).(Hasher[T])
	case int16:
		return (any)(int16Hasher{}).(Hasher[T])
	case int32:
		return (any)(int32Hasher{}).(Hasher[T])
	case int64:
		return (any)(int64Hasher{}).(Hasher[T])

	case uint:
		return (any)(uintHasher{}).(Hasher[T])
	case uint8:
		return (any)(uint8Hasher{}).(Hasher[T])
	case uint16:
		return (any)(uint16Hasher{}).(Hasher[T])
	case uint32:
		return (any)(uint32Hasher{}).(Hasher[T])
	case uint64:
		return (any)(uint64Hasher{}).(Hasher[T])

	case float32:
		return (any)(float32Hasher{}).(Hasher[T])
	case float64:
		return (any)(float64Hasher{}).(Hasher[T])

	case bin.Bin64:
		return (any)(bin64Hasher{}).(Hasher[T])
	case bin.Bin128:
		return (any)(bin128Hasher{}).(Hasher[T])
	case bin.Bin192:
		return (any)(bin192Hasher{}).(Hasher[T])
	case bin.Bin256:
		return (any)(bin256Hasher{}).(Hasher[T])

	case []byte:
		return (any)(bytesHasher{}).(Hasher[T])
	case string:
		return (any)(stringHasher{}).(Hasher[T])

	case Hashable:
		return (any)(hashableHasher[T]{}).(Hasher[T])
	}

	panic(fmt.Sprintf("unsupported hasher type: %T", k))
}

// internal

type (
	boolHasher struct{}

	intHasher   struct{}
	int8Hasher  struct{}
	int16Hasher struct{}
	int32Hasher struct{}
	int64Hasher struct{}

	uintHasher   struct{}
	uint8Hasher  struct{}
	uint16Hasher struct{}
	uint32Hasher struct{}
	uint64Hasher struct{}

	float32Hasher struct{}
	float64Hasher struct{}

	bin64Hasher  struct{}
	bin128Hasher struct{}
	bin192Hasher struct{}
	bin256Hasher struct{}

	bytesHasher  struct{}
	stringHasher struct{}

	hashableHasher[T any] struct{}
)

func (boolHasher) Hash32(v bool) uint32 { return HashBool(v) }

func (intHasher) Hash32(v int) uint32     { return HashInt(v) }
func (int8Hasher) Hash32(v int8) uint32   { return HashInt8(v) }
func (int16Hasher) Hash32(v int16) uint32 { return HashInt16(v) }
func (int32Hasher) Hash32(v int32) uint32 { return HashInt32(v) }
func (int64Hasher) Hash32(v int64) uint32 { return HashInt64(v) }

func (uintHasher) Hash32(v uint) uint32     { return HashUint(v) }
func (uint8Hasher) Hash32(v uint8) uint32   { return HashUint8(v) }
func (uint16Hasher) Hash32(v uint16) uint32 { return HashUint16(v) }
func (uint32Hasher) Hash32(v uint32) uint32 { return HashUint32(v) }
func (uint64Hasher) Hash32(v uint64) uint32 { return HashUint64(v) }

func (float32Hasher) Hash32(v float32) uint32 { return HashFloat32(v) }
func (float64Hasher) Hash32(v float64) uint32 { return HashFloat64(v) }

func (bin64Hasher) Hash32(v bin.Bin64) uint32   { return HashBin64(v) }
func (bin128Hasher) Hash32(v bin.Bin128) uint32 { return HashBin128(v) }
func (bin192Hasher) Hash32(v bin.Bin192) uint32 { return HashBin192(v) }
func (bin256Hasher) Hash32(v bin.Bin256) uint32 { return HashBin256(v) }

func (bytesHasher) Hash32(v []byte) uint32  { return HashBytes(v) }
func (stringHasher) Hash32(v string) uint32 { return HashString(v) }

func (hashableHasher[T]) Hash32(v T) uint32 { return (any)(v).(Hashable).Hash32() }
