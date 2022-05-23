package u128

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

// U128 is a 16 byte random or time-random unique value.
type U128 [ByteLen]byte

// FromInt64 converts an int64 into a big endian U128.
func FromInt64(v int64) U128 {
	u := U128{}
	b := u[8:]

	binary.BigEndian.PutUint64(b, uint64(v))
	return u
}

// Compare compares two IDs.
func (u0 U128) Compare(u1 U128) int {
	return bytes.Compare(u0[:], u1[:])
}

// IsZero returns if the id is zero.
func (u U128) IsZero() bool {
	return u == U128{}
}

// Less returns whether the current ID is less than another.
func (u0 U128) Less(u1 U128) bool {
	return bytes.Compare(u0[:], u1[:]) < 0
}

// Size returns 16 bytes.
func (u U128) Size() int {
	return ByteLen
}

// String returns a 32-char lower-case hex-encoded string.
func (u U128) String() string {
	b := make([]byte, CharLen)
	hex.Encode(b, u[:])
	return string(b)
}

// Marshal marshals the ID to a 16-byte array.
func (u U128) Marshal() ([]byte, error) {
	return u[:], nil
}

// MarshalTo marshals the ID to a 16-byte array.
func (u U128) MarshalTo(b []byte) (n int, err error) {
	copy(b, u[:])
	return len(u), nil
}

// Unmarshal parses a 16-byte array.
func (u *U128) Unmarshal(b []byte) error {
	u0, err := Parse(b)
	if err != nil {
		return err
	}

	*u = u0
	return nil
}

// MarshalJSON marshals the ID to a JSON string.
func (u U128) MarshalJSON() ([]byte, error) {
	s := u.String()
	return json.Marshal(s)
}

// MarshalJSON parses the ID from a JSON string.
func (u *U128) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	u0, err := ParseString(s)
	if err != nil {
		return err
	}

	*u = u0
	return nil
}
