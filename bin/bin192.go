// Copyright 2026 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package bin

import (
	"encoding/hex"
	"encoding/json"
	"regexp"
	"slices"

	"github.com/basecomplextech/baselibrary/buffer"
)

const (
	// Len192 is the length of a bin192 in bytes.
	Len192 = 24

	// Len192Char is the length of a bin192 encoded as a hex string with dashes.
	Len192Char = (Len192 * 2) + 2 // 2 for dashes
)

var (
	// Max192 is the maximum value of a bin192.
	Max192 = Bin192{Max64, Max64, Max64}

	// Regexp192 is a regular expression for validating a hex-encoded bin192 string.
	Regexp192 = regexp.MustCompile(`^[0-9A-Za-z]{16}-[0-9A-Za-z]{16}-[0-9A-Za-z]{16}$`)
)

// Bin192 is 24-byte binary value, represented as three bin64s.
type Bin192 [3]Bin64

// New192 returns a bin192 from a byte array.
func New192(p [Len192]byte) Bin192 {
	b := Bin192{}
	copy(b[0][:], p[:8])
	copy(b[1][:], p[8:16])
	copy(b[2][:], p[16:24])
	return b
}

// Int192 returns a bin192 from three int64s encoded as big-endian.
func Int192(v0, v1, v2 int64) Bin192 {
	return Bin192{
		Int64(v0),
		Int64(v1),
		Int64(v2),
	}
}

// Size returns 24 bytes.
func (b Bin192) Size() int {
	return Len192
}

// Compare

// Compare compares two values.
//
// Returns:
//
//	-1 if a < b
//	 0 if a == b
//	 1 if a > b
func (b Bin192) Compare(b1 Bin192) int {
	c := b[0].Compare(b1[0])
	if c != 0 {
		return c
	}
	c = b[1].Compare(b1[1])
	if c != 0 {
		return c
	}
	return b[2].Compare(b1[2])
}

// Equal returns whether two values are equal.
func (b Bin192) Equal(b1 Bin192) bool {
	return b == b1
}

// Less returns whether the current value is less than another.
func (b Bin192) Less(b1 Bin192) bool {
	return b.Compare(b1) < 0
}

// IsZero returns whether the value is zero.
func (b Bin192) IsZero() bool {
	return b == Bin192{}
}

// Ints

// Int64 returns three int64 decoded as big-endian.
func (b Bin192) Int64() [3]int64 {
	v0 := b[0].Int64()
	v1 := b[1].Int64()
	v2 := b[2].Int64()
	return [3]int64{v0, v1, v2}
}

// Uint64 returns three uint64 decoded as big-endian.
func (b Bin192) Uint64() [3]uint64 {
	v0 := b[0].Uint64()
	v1 := b[1].Uint64()
	v2 := b[2].Uint64()
	return [3]uint64{v0, v1, v2}
}

// Hash

// Hash32 returns a 32-bit hash.
// The method decodes the value as three big-endian uint64s and then xors their halves.
func (b Bin192) Hash32() uint32 {
	v0 := b[0].Hash64()
	v1 := b[1].Hash64()
	v2 := b[2].Hash64()

	v0 = v0 ^ (v0 >> 32)
	v1 = v1 ^ (v1 >> 32)
	v2 = v2 ^ (v2 >> 32)

	v := v0 ^ v1 ^ v2
	return uint32(v)
}

// Hash64 returns a 64-bit hash.
// The method decodes the value as three big-endian uint64s and then xors them.
func (b Bin192) Hash64() uint64 {
	v0 := b[0].Hash64()
	v1 := b[1].Hash64()
	v2 := b[2].Hash64()

	return v0 ^ v1 ^ v2
}

// String/Hex

// String returns a lower-case hex-encoded string.
func (b Bin192) String() string {
	buf := make([]byte, 0, Len192Char)
	buf = b.AppendHexTo(buf)
	return string(buf)
}

// AppendHexTo appends a lower-case hex-encoded string to a buffer.
func (b Bin192) AppendHexTo(buf []byte) []byte {
	p := [Len192]byte{}
	b.MarshalTo(p[:])

	n := len(buf)
	n1 := n + Len192Char

	buf = slices.Grow(buf, n1)
	buf = buf[:n1]

	hex.Encode(buf[n:], p[:8])
	buf[n+16] = '-'
	hex.Encode(buf[n+17:], p[8:16])
	buf[n+33] = '-'
	hex.Encode(buf[n+34:], p[16:24])
	return buf
}

// Marshal

// Marshal marshals the value to a 24-byte array.
func (b Bin192) Marshal() []byte {
	buf := make([]byte, Len192)
	b.MarshalTo(buf)
	return buf
}

// MarshalTo marshals the value to a 24-byte array.
func (b Bin192) MarshalTo(buf []byte) {
	b[0].MarshalTo(buf[:8])
	b[1].MarshalTo(buf[8:16])
	b[2].MarshalTo(buf[16:24])
}

// MarshalToBuffer marshals the value to a buffer.
func (b Bin192) MarshalToBuffer(buf buffer.Buffer) []byte {
	p := buf.Grow(Len192)
	b.MarshalTo(p)
	return p
}

// MarshalJSON marshals the value to a JSON string.
func (b Bin192) MarshalJSON() ([]byte, error) {
	s := b.String()
	return json.Marshal(s)
}

// Unmarshal

// Unmarshal parses a 24-byte array.
func (b *Bin192) Unmarshal(buf []byte) error {
	b0, err := Parse192(buf)
	if err != nil {
		return err
	}

	*b = b0
	return nil
}

// MarshalJSON parses the value from a JSON string.
func (b *Bin192) UnmarshalJSON(buf []byte) error {
	var s string
	if err := json.Unmarshal(buf, &s); err != nil {
		return err
	}

	b0, err := ParseString192(s)
	if err != nil {
		return err
	}

	*b = b0
	return nil
}
