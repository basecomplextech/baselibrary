package u128

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
)

const (
	byteLen     = 16
	byteTimeLen = 8
	charLen     = (byteLen * 2) + 1 // 341a7d60bc5893a6-4bda3de06721534c
)

// U128 is a 16 byte random or time-random unique value.
type U128 [byteLen]byte

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
	return byteLen
}

// String returns a 33-char lower-case hex-encoded string.
func (u U128) String() string {
	b := make([]byte, charLen)
	hex.Encode(b[:16], u[:8])
	b[16] = '-'
	hex.Encode(b[17:], u[8:])
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
