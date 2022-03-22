package u256

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
)

const (
	byteLen     = 32
	byteTimeLen = 8
	charLen     = (byteLen * 2) + 2 // 46d706492a2cd221-8b33e1dd3ab9bc98-a52898502c59bbb151644322f638df47
)

// U256 is a 32 byte value.
type U256 [byteLen]byte

// FromInt64 converts an int64 into a big endian U256.
func FromInt64(v int64) U256 {
	u := U256{}
	b := u[24:]

	binary.BigEndian.PutUint64(b, uint64(v))
	return u
}

// Bytes marshals an id to bytes.
func (u U256) Bytes() []byte {
	b := make([]byte, len(u))
	copy(b, u[:])
	return b
}

// Compare compares two IDs.
func (u0 U256) Compare(u1 U256) int {
	return bytes.Compare(u0[:], u1[:])
}

// Equal returns whether two IDs are equal.
func (u0 U256) Equal(u1 U256) bool {
	return u0 == u1
}

// IsZero returns true if the id is zero.
func (u U256) IsZero() bool {
	return u == U256{}
}

// Less returns whether the current ID is less than another.
func (u0 U256) Less(u1 U256) bool {
	return bytes.Compare(u0[:], u1[:]) < 0
}

// Size returns 32 bytes.
func (u U256) Size() int {
	return byteLen
}

// String returns a 32-char lower-case hex-encoded string.
func (u U256) String() string {
	b := make([]byte, charLen)

	// 46d706492a2cd221-8b33e1dd3ab9bc98-a52898502c59bbb151644322f638df47
	hex.Encode(b[:16], u[:8])
	b[16] = '-'
	hex.Encode(b[17:33], u[8:16])
	b[33] = '-'
	hex.Encode(b[34:], u[16:])
	return string(b)
}

// Marshal marshals the ID to a 32-byte array.
func (u U256) Marshal() ([]byte, error) {
	return u[:], nil
}

// MarshalTo marshals the ID to a 32-byte array.
func (u U256) MarshalTo(b []byte) (n int, err error) {
	copy(b, u[:])
	return len(u), nil
}

// Unmarshal parses a 32-byte array.
func (u *U256) Unmarshal(b []byte) error {
	u0, err := Parse(b)
	if err != nil {
		return err
	}

	*u = u0
	return nil
}

// MarshalJSON marshals the ID to a JSON string.
func (u U256) MarshalJSON() ([]byte, error) {
	s := u.String()
	return json.Marshal(s)
}

// MarshalJSON parses the ID from a JSON string.
func (u *U256) UnmarshalJSON(b []byte) error {
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
