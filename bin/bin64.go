package bin

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"regexp"
)

const (
	ByteLen64 = 8
	CharLen64 = (ByteLen64 * 2) // 341a7d60bc5893a6
)

var Pattern64 = regexp.MustCompile(`^[0-9A-Za-z]{16}$`)

// Bin64 is a binary 64-bit value.
type Bin64 [ByteLen64]byte

// Bin64FromInt converts an int into a bin64.
func Bin64FromInt(v int) Bin64 {
	b := Bin64{}
	binary.BigEndian.PutUint64(b[:], uint64(v))
	return b
}

// Compare compares two values.
func (b0 Bin64) Compare(b1 Bin64) int {
	return bytes.Compare(b0[:], b1[:])
}

// Less returns whether the current value is less than another.
func (b0 Bin64) Less(b1 Bin64) bool {
	return bytes.Compare(b0[:], b1[:]) < 0
}

// Size returns 8 bytes.
func (b Bin64) Size() int {
	return len(b)
}

// String returns a 16-char lower-case hex-encoded string.
func (b Bin64) String() string {
	buf := make([]byte, CharLen64)
	hex.Encode(buf, b[:])
	return string(buf)
}

// Marshal marshals the value to a 16-byte array.
func (b Bin64) Marshal() ([]byte, error) {
	return b[:], nil
}

// MarshalTo marshals the value to a 16-byte array.
func (b Bin64) MarshalTo(buf []byte) (n int, err error) {
	copy(buf, b[:])
	return len(b), nil
}

// Unmarshal parses a 16-byte array.
func (b *Bin64) Unmarshal(buf []byte) error {
	b0, err := Parse64(buf)
	if err != nil {
		return err
	}

	*b = b0
	return nil
}

// MarshalJSON marshals the value to a JSON string.
func (b Bin64) MarshalJSON() ([]byte, error) {
	s := b.String()
	return json.Marshal(s)
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
