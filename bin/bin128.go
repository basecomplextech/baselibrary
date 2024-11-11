// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package bin

import (
	"encoding/hex"
	"encoding/json"
	"regexp"

	"github.com/basecomplextech/baselibrary/buffer"
)

const Len128 = 16
const Len128Char = (Len128 * 2) + 1 // 341a7d60bc5893a6-4bda3de06721534c

var Max128 = Bin128{Max64, Max64}
var Regexp128 = regexp.MustCompile(`^[0-9A-Za-z]{16}-[0-9A-Za-z]{16}$`)

// Bin128

// Bin128 is a 128-bit value.
type Bin128 [2]Bin64

// New128 returns a bin128 from a byte array.
func New128(b [Len128]byte) Bin128 {
	v := Bin128{}
	copy(v[0][:], b[0:])
	copy(v[1][:], b[8:])
	return v
}

// Int128 returns a bin128 from two int64 encoded as big-endian.
func Int128(v0, v1 int64) Bin128 {
	return Bin128{
		Int64(v0),
		Int64(v1),
	}
}

// Size returns 16 bytes.
func (b Bin128) Size() int {
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
func (b0 Bin128) Compare(b1 Bin128) int {
	cmp := b0[0].Compare(b1[0])
	if cmp != 0 {
		return cmp
	}
	return b0[1].Compare(b1[1])
}

// Equal returns whether two values are equal.
func (b0 Bin128) Equal(b1 Bin128) bool {
	return b0 == b1
}

// Less returns whether the current value is less than another.
func (b0 Bin128) Less(b1 Bin128) bool {
	return b0.Compare(b1) < 0
}

// Ints

// Int64 returns two int64 decoded as big-endian.
func (b Bin128) Int64() (int64, int64) {
	v0 := b[0].Int64()
	v1 := b[1].Int64()
	return v0, v1
}

// Uint64 returns two uint64 decoded as big-endian.
func (b Bin128) Uint64() (uint64, uint64) {
	v0 := b[0].Uint64()
	v1 := b[1].Uint64()
	return v0, v1
}

// Hash

// Hash32 returns a 32-bit hash.
// The method decodes the value as two big-endian uint64s and then xors their halves.
func (b Bin128) Hash32() uint32 {
	v0 := b[0].Uint64()
	v1 := b[1].Uint64()

	v0 = v0 ^ (v0 >> 32)
	v1 = v1 ^ (v1 >> 32)

	v := v0 ^ v1
	return uint32(v)
}

// Hash64 returns a 64-bit hash.
// The method decodes the value as two big-endian uint64s and then xors them.
func (b Bin128) Hash64() uint64 {
	v0 := b[0].Uint64()
	v1 := b[1].Uint64()
	return v0 ^ v1
}

// String/Hex

// String returns a 33-char lower-case hex-encoded string.
func (b Bin128) String() string {
	buf := make([]byte, Len128Char)
	hex.Encode(buf, b[0][:])
	buf[16] = '-'
	hex.Encode(buf[17:], b[1][:])
	return string(buf)
}

// AppendHexTo appends a 33-char lower-case hex-encoded string to a buffer.
func (b Bin128) AppendHexTo(buf []byte) []byte {
	n := len(buf)
	n1 := n + Len128Char

	buf = buf[:n1]
	hex.Encode(buf[n:], b[0][:])
	buf[n+16] = '-'
	hex.Encode(buf[n+17:], b[1][:])
	return buf
}

// Marshal

// Marshal marshals the value to a 16-byte array.
func (b Bin128) Marshal() []byte {
	buf := make([]byte, Len128)
	b.MarshalTo(buf)
	return buf
}

// MarshalTo marshals the value to a 16-byte array.
func (b Bin128) MarshalTo(buf []byte) {
	copy(buf, b[0][:])
	copy(buf[8:], b[1][:])
}

// MarshalToBuffer marshals the value to a buffer.
func (b Bin128) MarshalToBuffer(buf buffer.Buffer) []byte {
	p := buf.Grow(Len128)
	b.MarshalTo(p)
	return p
}

// MarshalJSON marshals the value to a JSON string.
func (b Bin128) MarshalJSON() ([]byte, error) {
	s := b.String()
	return json.Marshal(s)
}

// Unmarshal

// Unmarshal parses a 16-byte array.
func (b *Bin128) Unmarshal(buf []byte) error {
	b0, err := Parse128(buf)
	if err != nil {
		return err
	}

	*b = b0
	return nil
}

// MarshalJSON parses the value from a JSON string.
func (b *Bin128) UnmarshalJSON(buf []byte) error {
	var s string
	if err := json.Unmarshal(buf, &s); err != nil {
		return err
	}

	b0, err := ParseString128(s)
	if err != nil {
		return err
	}

	*b = b0
	return nil
}
