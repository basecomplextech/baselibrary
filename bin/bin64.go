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

var (
	Max64     = Bin64{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	Pattern64 = regexp.MustCompile(`^[0-9A-Za-z]{16}$`)
)

// Bin64 is a binary 64-bit value.
type Bin64 [ByteLen64]byte

// Int64 returns a bin64 from an int64, encoded as big-endian.
func Int64(v int64) Bin64 {
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

// Int64 returns an int64 decoded as big-endian.
func (b Bin64) Int64() int64 {
	return int64(binary.BigEndian.Uint64(b[:]))
}

// Hash32 returns a 32-bit hash of the value.
// The method decodes the value as a big-endian uint64 and then xors the two halves.
func (b Bin64) Hash32() uint32 {
	v := binary.BigEndian.Uint64(b[:])
	v = v ^ (v >> 32) // xor of two halves
	return uint32(v)
}

// Hash64 returns a 64-bit hash of the value.
// The method decodes the value as a big-endian uint64.
func (b Bin64) Hash64() uint64 {
	return binary.BigEndian.Uint64(b[:])
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

// AppendHexTo appends a 16-char lower-case hex-encoded string to a buffer.
func (b Bin64) AppendHexTo(buf []byte) []byte {
	n := len(buf)
	n1 := n + CharLen64

	buf = buf[:n1]
	hex.Encode(buf[n:], b[:])
	return buf
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
