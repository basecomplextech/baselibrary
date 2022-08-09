// Package bin128 provides a binary 128-bit value.
package bin128

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
)

const (
	ByteLen     = 16
	CharLen     = (ByteLen * 2) // 341a7d60bc5893a64bda3de06721534c
	byteTimeLen = 8
)

// B128 is a 16 byte random or time-random unique value.
type B128 [ByteLen]byte

// FromInt64 converts an int64 into a big endian B128.
func FromInt64(v int64) B128 {
	b := B128{}
	buf := b[8:]

	binary.BigEndian.PutUint64(buf, uint64(v))
	return b
}

// Compare compares two IDs.
func (b0 B128) Compare(b1 B128) int {
	return bytes.Compare(b0[:], b1[:])
}

// IsZero returns if the id is zero.
func (b B128) IsZero() bool {
	return b == B128{}
}

// Less returns whether the current ID is less than another.
func (b0 B128) Less(b1 B128) bool {
	return bytes.Compare(b0[:], b1[:]) < 0
}

// Size returns 16 bytes.
func (b B128) Size() int {
	return ByteLen
}

// String returns a 32-char lower-case hex-encoded string.
func (b B128) String() string {
	buf := make([]byte, CharLen)
	hex.Encode(buf, b[:])
	return string(buf)
}

// Marshal marshals the ID to a 16-byte array.
func (b B128) Marshal() ([]byte, error) {
	return b[:], nil
}

// MarshalTo marshals the ID to a 16-byte array.
func (b B128) MarshalTo(buf []byte) (n int, err error) {
	copy(buf, b[:])
	return len(b), nil
}

// Unmarshal parses a 16-byte array.
func (b *B128) Unmarshal(buf []byte) error {
	b0, err := Parse(buf)
	if err != nil {
		return err
	}

	*b = b0
	return nil
}

// MarshalJSON marshals the ID to a JSON string.
func (b B128) MarshalJSON() ([]byte, error) {
	s := b.String()
	return json.Marshal(s)
}

// MarshalJSON parses the ID from a JSON string.
func (b *B128) UnmarshalJSON(buf []byte) error {
	var s string
	if err := json.Unmarshal(buf, &s); err != nil {
		return err
	}

	b0, err := ParseString(s)
	if err != nil {
		return err
	}

	*b = b0
	return nil
}
