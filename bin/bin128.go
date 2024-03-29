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

var (
	Max128 = Bin128{
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	}
	Pattern128 = regexp.MustCompile(`^[0-9A-Za-z]{32}$`)
)

// Bin128 is a 128-bit value.
type Bin128 [ByteLen128]byte

// Int128 returns a bin128 from two int64 encoded as big-endian.
func Int128(v0, v1 int64) Bin128 {
	b := Bin128{}
	binary.BigEndian.PutUint64(b[:], uint64(v0))
	binary.BigEndian.PutUint64(b[8:], uint64(v1))
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

// Int128 returns two int64s decoded as big-endian.
func (b Bin128) Int128() (int64, int64) {
	v0 := binary.BigEndian.Uint64(b[:])
	v1 := binary.BigEndian.Uint64(b[8:])
	return int64(v0), int64(v1)
}

// Hash32 returns a 32-bit hash.
// The method decodes the value as two big-endian uint64s and then xors their halves.
func (b Bin128) Hash32() uint32 {
	v0 := binary.BigEndian.Uint64(b[:])
	v1 := binary.BigEndian.Uint64(b[8:])

	v0 = v0 ^ (v0 >> 32)
	v1 = v1 ^ (v1 >> 32)

	v := v0 ^ v1
	return uint32(v)
}

// Hash64 returns a 64-bit hash.
// The method decodes the value as two big-endian uint64s and then xors them.
func (b Bin128) Hash64() uint64 {
	v0 := binary.BigEndian.Uint64(b[:])
	v1 := binary.BigEndian.Uint64(b[8:])
	return v0 ^ v1
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

// AppendHexTo appends a 32-char lower-case hex-encoded string to a buffer.
func (b Bin128) AppendHexTo(buf []byte) []byte {
	n := len(buf)
	n1 := n + CharLen128

	buf = buf[:n1]
	hex.Encode(buf[n:], b[:])
	return buf
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
