// Package bin256 provides a binary 256-bit value.
package types

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
)

const (
	Bin256ByteLen = 32
	Bin256CharLen = (Bin256ByteLen * 2)
)

// Bin256 is a 32 byte value.
type Bin256 [Bin256ByteLen]byte

// Bin256FromInt64 converts an int64 into a big endian Bin256.
func Bin256FromInt64(v int64) Bin256 {
	b := Bin256{}
	buf := b[24:]

	binary.BigEndian.PutUint64(buf, uint64(v))
	return b
}

// Bytes marshals an id to bytes.
func (b Bin256) Bytes() []byte {
	buf := make([]byte, len(b))
	copy(buf, b[:])
	return buf
}

// Compare compares two values.
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

// String returns a 32-char lower-case hex-encoded string.
func (b Bin256) String() string {
	buf := make([]byte, Bin256CharLen)
	hex.Encode(buf, b[:])
	return string(buf)
}

// Marshal marshals the value to a 32-byte array.
func (b Bin256) Marshal() ([]byte, error) {
	return b[:], nil
}

// MarshalTo marshals the value to a 32-byte array.
func (b Bin256) MarshalTo(buf []byte) (n int, err error) {
	copy(buf, b[:])
	return len(b), nil
}

// Unmarshal parses a 32-byte array.
func (b *Bin256) Unmarshal(buf []byte) error {
	b0, err := ParseBin256(buf)
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

	b0, err := ParseStringBin256(s)
	if err != nil {
		return err
	}

	*b = b0
	return nil
}
