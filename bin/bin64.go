// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package bin

import (
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"regexp"
	"slices"

	"github.com/basecomplextech/baselibrary/buffer"
)

const Len64 = 8
const Len64Char = (Len64 * 2) // 341a7d60bc5893a6

var Max64 = Bin64{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
var Regexp64 = regexp.MustCompile(`^[0-9A-Za-z]{16}$`)

// Bin64

// Bin64 is a binary 64-bit value.
type Bin64 [8]byte

// New64 returns a bin64 from a byte array.
func New64(b [Len64]byte) Bin64 {
	return Bin64(b)
}

// Int64 returns a bin64 from an int64, encoded as big-endian.
func Int64(v int64) Bin64 {
	b := Bin64{}
	binary.BigEndian.PutUint64(b[:], uint64(v))
	return b
}

// Size returns 8 bytes.
func (b Bin64) Size() int {
	return Len64
}

// Compare

// Compare compares two values.
//
// Returns:
//
//	-1 if a < b
//	 0 if a == b
//	 1 if a > b
func (b Bin64) Compare(b1 Bin64) int {
	i := b.Uint64()
	i1 := b1.Uint64()

	switch {
	case i == i1:
		return 0
	case i < i1:
		return -1
	default:
		return 1
	}
}

// Equal returns whether two values are equal.
func (b Bin64) Equal(b1 Bin64) bool {
	return b == b1
}

// Less returns whether the current value is less than another.
func (b Bin64) Less(b1 Bin64) bool {
	return b.Compare(b1) < 0
}

// IsZero returns whether the value is zero.
func (b Bin64) IsZero() bool {
	return b == Bin64{}
}

// Ints

// Int64 returns an int64 decoded as big-endian.
func (b Bin64) Int64() int64 {
	return int64(binary.BigEndian.Uint64(b[:]))
}

// Uint64 returns a uint64 decoded as big-endian.
func (b Bin64) Uint64() uint64 {
	return binary.BigEndian.Uint64(b[:])
}

// Hash

// Hash32 returns a 32-bit hash of the value.
// The method decodes the value as a big-endian uint64 and then xors the two halves.
func (b Bin64) Hash32() uint32 {
	i := b.Uint64()
	h := i ^ (i >> 32) // xor of two halves
	return uint32(h)
}

// Hash64 returns a 64-bit hash of the value.
// The method decodes the value as a big-endian uint64.
func (b Bin64) Hash64() uint64 {
	return b.Uint64()
}

// String/Hex

// String returns a 16-char lower-case hex-encoded string.
func (b Bin64) String() string {
	buf := make([]byte, Len64Char)
	hex.Encode(buf, b[:])
	return string(buf)
}

// AppendHexTo appends a 16-char lower-case hex-encoded string to a buffer.
func (b Bin64) AppendHexTo(buf []byte) []byte {
	n := len(buf)
	n1 := n + Len64Char

	buf = slices.Grow(buf, n1)
	buf = buf[:n1]

	hex.Encode(buf[n:], b[:])
	return buf
}

// Marshal

// Marshal marshals the value to a 16-byte array.
func (b Bin64) Marshal() []byte {
	return b[:]
}

// MarshalTo marshals the value to a 16-byte array.
func (b Bin64) MarshalTo(buf []byte) {
	copy(buf, b[:])
}

// MarshalToBuffer marshals the value to a buffer.
func (b Bin64) MarshalToBuffer(buf buffer.Buffer) []byte {
	p := buf.Grow(Len64)
	copy(p, b[:])
	return p
}

// MarshalJSON marshals the value to a JSON string.
func (b Bin64) MarshalJSON() ([]byte, error) {
	s := b.String()
	return json.Marshal(s)
}

// Unmarshal

// Unmarshal parses a 16-byte array.
func (b *Bin64) Unmarshal(buf []byte) error {
	b0, err := Parse64(buf)
	if err != nil {
		return err
	}

	*b = b0
	return nil
}

// MarshalJSON parses the value from a JSON string.
func (b *Bin64) UnmarshalJSON(buf []byte) error {
	var s string
	if err := json.Unmarshal(buf, &s); err != nil {
		return err
	}

	b0, err := ParseString64(s)
	if err != nil {
		return err
	}

	*b = b0
	return nil
}
