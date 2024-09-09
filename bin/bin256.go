// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

// Package bin256 provides a binary 256-bit value.
package bin

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"regexp"

	"github.com/basecomplextech/baselibrary/buffer"
)

const (
	ByteLen256 = 32
	CharLen256 = (ByteLen256 * 2) + 3
)

var (
	Max256 = Bin256{
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	}
	Pattern256 = regexp.MustCompile(`^[0-9A-Za-z]{16}-[0-9A-Za-z]{16}-[0-9A-Za-z]{16}-[0-9A-Za-z]{16}$`)
)

// Bin256 is a 32 byte value.
type Bin256 [ByteLen256]byte

// Int256 returns a bin256 from four int64s encoded as big-endian.
func Int256(v0, v1, v2, v3 int64) Bin256 {
	b := Bin256{}
	binary.BigEndian.PutUint64(b[:], uint64(v0))
	binary.BigEndian.PutUint64(b[8:], uint64(v1))
	binary.BigEndian.PutUint64(b[16:], uint64(v2))
	binary.BigEndian.PutUint64(b[24:], uint64(v3))
	return b
}

// Join256 joins four bin64 values into a bin256.
func Join256(b0, b1, b2, b3 Bin64) Bin256 {
	b := Bin256{}
	copy(b[:], b0[:])
	copy(b[8:], b1[:])
	copy(b[16:], b2[:])
	copy(b[24:], b3[:])
	return b
}

// Bytes marshals an id to bytes.
func (b Bin256) Bytes() []byte {
	buf := make([]byte, len(b))
	copy(buf, b[:])
	return buf
}

// Compare compares two values.
//
// Returns:
//
//	-1 if a < b
//	 0 if a == b
//	 1 if a > b
func (b0 Bin256) Compare(b1 Bin256) int {
	return bytes.Compare(b0[:], b1[:])
}

// Equal returns whether two values are equal.
func (b0 Bin256) Equal(b1 Bin256) bool {
	return b0 == b1
}

// Less returns whether the current value is less than another.
func (b0 Bin256) Less(b1 Bin256) bool {
	return bytes.Compare(b0[:], b1[:]) < 0
}

// Size returns 32 bytes.
func (b Bin256) Size() int {
	return len(b)
}

// Parts

// Int64s returns four int64s decoded as big-endian.
func (b Bin256) Int64s() [4]int64 {
	v0 := binary.BigEndian.Uint64(b[:])
	v1 := binary.BigEndian.Uint64(b[8:])
	v2 := binary.BigEndian.Uint64(b[16:])
	v3 := binary.BigEndian.Uint64(b[24:])
	return [4]int64{int64(v0), int64(v1), int64(v2), int64(v3)}
}

// Bin64s returns four bin64 values.
func (b Bin256) Bin64s() [4]Bin64 {
	var b0, b1, b2, b3 Bin64
	copy(b0[:], b[:])
	copy(b1[:], b[8:])
	copy(b2[:], b[16:])
	copy(b3[:], b[24:])
	return [4]Bin64{b0, b1, b2, b3}
}

// Hash

// Hash32 returns a 32-bit hash.
// The method decodes the value as four big-endian uint64s and then xors their halves.
func (b Bin256) Hash32() uint32 {
	v0 := binary.BigEndian.Uint64(b[:])
	v1 := binary.BigEndian.Uint64(b[8:])
	v2 := binary.BigEndian.Uint64(b[16:])
	v3 := binary.BigEndian.Uint64(b[24:])

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
	v0 := binary.BigEndian.Uint64(b[:])
	v1 := binary.BigEndian.Uint64(b[8:])
	v2 := binary.BigEndian.Uint64(b[16:])
	v3 := binary.BigEndian.Uint64(b[24:])

	v := v0 ^ v1 ^ v2 ^ v3
	return v
}

// String/Hex

// String returns a 64-char lower-case hex-encoded string.
func (b Bin256) String() string {
	buf := make([]byte, CharLen256)
	hex.Encode(buf, b[:8])
	buf[16] = '-'
	hex.Encode(buf[17:], b[8:16])
	buf[33] = '-'
	hex.Encode(buf[34:], b[16:24])
	buf[50] = '-'
	hex.Encode(buf[51:], b[24:])
	return string(buf)
}

// AppendHexTo appends a 67-char lower-case hex-encoded string to a buffer.
func (b Bin256) AppendHexTo(buf []byte) []byte {
	n := len(buf)
	n1 := n + CharLen256

	buf = buf[:n1]
	hex.Encode(buf[n:], b[:8])
	buf[n+16] = '-'
	hex.Encode(buf[n+17:], b[8:16])
	buf[n+33] = '-'
	hex.Encode(buf[n+34:], b[16:24])
	buf[n+50] = '-'
	hex.Encode(buf[n+51:], b[24:])
	return buf
}

// Marshal

// Marshal marshals the value to a 32-byte array.
func (b Bin256) Marshal() []byte {
	return b[:]
}

// MarshalTo marshals the value to a 32-byte array.
func (b Bin256) MarshalTo(buf []byte) {
	copy(buf, b[:])
}

// MarshalToBuffer marshals the value to a buffer.
func (b Bin256) MarshalToBuffer(buf buffer.Buffer) []byte {
	p := buf.Grow(ByteLen256)
	copy(p, b[:])
	return p
}

// Unmarshal parses a 32-byte array.
func (b *Bin256) Unmarshal(buf []byte) error {
	b0, err := Parse256(buf)
	if err != nil {
		return err
	}

	*b = b0
	return nil
}

// MarshalJSON marshals the value to a JSON string.
func (b Bin256) MarshalJSON() ([]byte, error) {
	s := b.String()
	return json.Marshal(s)
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
