// Package bin256 provides a binary 256-bit value.
package bin256

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
)

const (
	ByteLen = 32
	CharLen = (ByteLen * 2)

	byteTimeLen = 8
)

// B256 is a 32 byte value.
type B256 [ByteLen]byte

// FromInt64 converts an int64 into a big endian B256.
func FromInt64(v int64) B256 {
	b := B256{}
	buf := b[24:]

	binary.BigEndian.PutUint64(buf, uint64(v))
	return b
}

// Bytes marshals an id to bytes.
func (b B256) Bytes() []byte {
	buf := make([]byte, len(b))
	copy(buf, b[:])
	return buf
}

// Compare compares two IDs.
func (b0 B256) Compare(b1 B256) int {
	return bytes.Compare(b0[:], b1[:])
}

// Equal returns whether two IDs are equal.
func (b0 B256) Equal(b1 B256) bool {
	return b0 == b1
}

// IsZero returns true if the id is zero.
func (b B256) IsZero() bool {
	return b == B256{}
}

// Less returns whether the current ID is less than another.
func (b0 B256) Less(b1 B256) bool {
	return bytes.Compare(b0[:], b1[:]) < 0
}

// Size returns 32 bytes.
func (b B256) Size() int {
	return ByteLen
}

// String returns a 32-char lower-case hex-encoded string.
func (b B256) String() string {
	buf := make([]byte, CharLen)
	hex.Encode(buf, b[:])
	return string(buf)
}

// Marshal marshals the ID to a 32-byte array.
func (b B256) Marshal() ([]byte, error) {
	return b[:], nil
}

// MarshalTo marshals the ID to a 32-byte array.
func (b B256) MarshalTo(buf []byte) (n int, err error) {
	copy(buf, b[:])
	return len(b), nil
}

// Unmarshal parses a 32-byte array.
func (b *B256) Unmarshal(buf []byte) error {
	b0, err := Parse(buf)
	if err != nil {
		return err
	}

	*b = b0
	return nil
}

// MarshalJSON marshals the ID to a JSON string.
func (b B256) MarshalJSON() ([]byte, error) {
	s := b.String()
	return json.Marshal(s)
}

// MarshalJSON parses the ID from a JSON string.
func (b *B256) UnmarshalJSON(buf []byte) error {
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
