// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

// Package bin256 provides a binary 256-bit value.
package bin

import (
	"encoding/hex"
	"encoding/json"
	"regexp"
	"slices"

	"github.com/basecomplextech/baselibrary/buffer"
)

const Len256 = 32
const Len256Char = (Len256 * 2) + 3

var Max256 = Bin256{Max64, Max64, Max64, Max64}
var Regexp256 = regexp.MustCompile(`^[0-9A-Za-z]{16}-[0-9A-Za-z]{16}-[0-9A-Za-z]{16}-[0-9A-Za-z]{16}$`)

// Bin256

// Bin256 is a 32 byte value.
type Bin256 [4]Bin64

// New256 returns a bin256 from a byte array.
func New256(p [Len256]byte) Bin256 {
	b := Bin256{}
	copy(b[0][:], p[:8])
	copy(b[1][:], p[8:16])
	copy(b[2][:], p[16:24])
	copy(b[3][:], p[24:])
	return b
}

// Int256 returns a bin256 from four int64s encoded as big-endian.
func Int256(v0, v1, v2, v3 int64) Bin256 {
	return Bin256{
		Int64(v0),
		Int64(v1),
		Int64(v2),
		Int64(v3),
	}
}

// Size returns 32 bytes.
func (b Bin256) Size() int {
	return Len256
}

// Compare

// Compare compares two values.
//
// Returns:
//
//	-1 if a < b
//	 0 if a == b
//	 1 if a > b
func (b Bin256) Compare(b1 Bin256) int {
	c := b[0].Compare(b1[0])
	if c != 0 {
		return c
	}
	c = b[1].Compare(b1[1])
	if c != 0 {
		return c
	}
	c = b[2].Compare(b1[2])
	if c != 0 {
		return c
	}
	return b[3].Compare(b1[3])
}

// Equal returns whether two values are equal.
func (b Bin256) Equal(b1 Bin256) bool {
	return b == b1
}

// Less returns whether the current value is less than another.
func (b Bin256) Less(b1 Bin256) bool {
	return b.Compare(b1) < 0
}

// IsZero returns whether the value is zero.
func (b Bin256) IsZero() bool {
	return b == Bin256{}
}

// Ints

// Int64 returns four int64 decoded as big-endian.
func (b Bin256) Int64() [4]int64 {
	v0 := b[0].Int64()
	v1 := b[1].Int64()
	v2 := b[2].Int64()
	v3 := b[3].Int64()
	return [4]int64{v0, v1, v2, v3}
}

// Uint64 returns four uint64 decoded as big-endian.
func (b Bin256) Uint64() [4]uint64 {
	v0 := b[0].Uint64()
	v1 := b[1].Uint64()
	v2 := b[2].Uint64()
	v3 := b[3].Uint64()
	return [4]uint64{v0, v1, v2, v3}
}

// Bin128 returns two 128-bit values.
func (b Bin256) Bin128() [2]Bin128 {
	return [2]Bin128{
		{b[0], b[1]},
		{b[2], b[3]},
	}
}

// Hash

// Hash32 returns a 32-bit hash.
// The method decodes the value as four big-endian uint64s and then xors their halves.
func (b Bin256) Hash32() uint32 {
	v0 := b[0].Hash64()
	v1 := b[1].Hash64()
	v2 := b[2].Hash64()
	v3 := b[3].Hash64()

	v0 = v0 ^ (v0 >> 32)
	v1 = v1 ^ (v1 >> 32)
	v2 = v2 ^ (v2 >> 32)
	v3 = v3 ^ (v3 >> 32)

	v := v0 ^ v1 ^ v2 ^ v3
	return uint32(v)
}

// Hash64 returns a 64-bit hash.
// The method decodes the value as four big-endian uint64s and then xors them.
func (b Bin256) Hash64() uint64 {
	v0 := b[0].Hash64()
	v1 := b[1].Hash64()
	v2 := b[2].Hash64()
	v3 := b[3].Hash64()

	v := v0 ^ v1 ^ v2 ^ v3
	return v
}

// String/Hex

// String returns a 64-char lower-case hex-encoded string.
func (b Bin256) String() string {
	buf := make([]byte, 0, Len256)
	buf = b.AppendHexTo(buf)
	return string(buf)
}

// AppendHexTo appends a 67-char lower-case hex-encoded string to a buffer.
func (b Bin256) AppendHexTo(buf []byte) []byte {
	p := [Len256]byte{}
	b.MarshalTo(p[:])

	n := len(buf)
	n1 := n + Len256Char

	buf = slices.Grow(buf, n1)
	buf = buf[:n1]

	hex.Encode(buf[n:], p[:8])
	buf[n+16] = '-'
	hex.Encode(buf[n+17:], p[8:16])
	buf[n+33] = '-'
	hex.Encode(buf[n+34:], p[16:24])
	buf[n+50] = '-'
	hex.Encode(buf[n+51:], p[24:])
	return buf
}

// Marshal

// Marshal marshals the value to a 32-byte array.
func (b Bin256) Marshal() []byte {
	buf := make([]byte, Len256)
	b.MarshalTo(buf)
	return buf
}

// MarshalTo marshals the value to a 32-byte array.
func (b Bin256) MarshalTo(buf []byte) {
	b[0].MarshalTo(buf[:8])
	b[1].MarshalTo(buf[8:16])
	b[2].MarshalTo(buf[16:24])
	b[3].MarshalTo(buf[24:])
}

// MarshalToBuffer marshals the value to a buffer.
func (b Bin256) MarshalToBuffer(buf buffer.Buffer) []byte {
	p := buf.Grow(Len256)
	b.MarshalTo(p)
	return p
}

// MarshalJSON marshals the value to a JSON string.
func (b Bin256) MarshalJSON() ([]byte, error) {
	s := b.String()
	return json.Marshal(s)
}

// Unmarshal

// Unmarshal parses a 32-byte array.
func (b *Bin256) Unmarshal(buf []byte) error {
	b0, err := Parse256(buf)
	if err != nil {
		return err
	}

	*b = b0
	return nil
}

// MarshalJSON parses the value from a JSON string.
func (b *Bin256) UnmarshalJSON(buf []byte) error {
	var s string
	if err := json.Unmarshal(buf, &s); err != nil {
		return err
	}

	b0, err := ParseString256(s)
	if err != nil {
		return err
	}

	*b = b0
	return nil
}
