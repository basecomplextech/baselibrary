package bin

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"regexp"
)

const (
	ByteLen128 = 16
	CharLen128 = (ByteLen128 * 2) // 341a7d60bc5893a64bda3de06721534c
)

var Pattern128 = regexp.MustCompile(`^[0-9A-Za-z]{32}$`)

// Bin128 is a 128-bit value.
type Bin128 [ByteLen128]byte

// Bin128FromInt converts an int into a bin128.
func Bin128FromInt(v int) Bin128 {
	b := Bin128{}
	buf := b[8:]

	binary.BigEndian.PutUint64(buf, uint64(v))
	return b
}

// Compare compares two values.
func (b0 Bin128) Compare(b1 Bin128) int {
	return bytes.Compare(b0[:], b1[:])
}

// Less returns whether the current value is less than another.
func (b0 Bin128) Less(b1 Bin128) bool {
	return bytes.Compare(b0[:], b1[:]) < 0
}

// Size returns 16 bytes.
func (b Bin128) Size() int {
	return len(b)
}

// String returns a 32-char lower-case hex-encoded string.
func (b Bin128) String() string {
	buf := make([]byte, CharLen128)
	hex.Encode(buf, b[:])
	return string(buf)
}

// Marshal marshals the value to a 16-byte array.
func (b Bin128) Marshal() ([]byte, error) {
	return b[:], nil
}

// MarshalTo marshals the value to a 16-byte array.
func (b Bin128) MarshalTo(buf []byte) (n int, err error) {
	copy(buf, b[:])
	return len(b), nil
}

// Unmarshal parses a 16-byte array.
func (b *Bin128) Unmarshal(buf []byte) error {
	b0, err := Parse128(buf)
	if err != nil {
		return err
	}

	*b = b0
	return nil
}

// MarshalJSON marshals the value to a JSON string.
func (b Bin128) MarshalJSON() ([]byte, error) {
	s := b.String()
	return json.Marshal(s)
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