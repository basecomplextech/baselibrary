// Package bin256 provides a binary 256-bit value.
package bin

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"regexp"
)

const (
	ByteLen256 = 32
	CharLen256 = (ByteLen256 * 2)
)

var (
	Max256 = Bin256{
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	}
	Pattern256 = regexp.MustCompile(`^[0-9A-Za-z]{64}$`)
)

// Bin256 is a 32 byte value.
type Bin256 [ByteLen256]byte

// Int256 returns a Bin256 from four int64s encoded as big-endian.
func Int256(v0, v1, v2, v3 int64) Bin256 {
	b := Bin256{}
	binary.BigEndian.PutUint64(b[:], uint64(v0))
	binary.BigEndian.PutUint64(b[8:], uint64(v1))
	binary.BigEndian.PutUint64(b[16:], uint64(v2))
	binary.BigEndian.PutUint64(b[24:], uint64(v3))
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

// Int256 returns four int64s decoded as big-endian.
func (b Bin256) Int256() [4]int64 {
	v0 := binary.BigEndian.Uint64(b[:])
	v1 := binary.BigEndian.Uint64(b[8:])
	v2 := binary.BigEndian.Uint64(b[16:])
	v3 := binary.BigEndian.Uint64(b[24:])
	return [4]int64{int64(v0), int64(v1), int64(v2), int64(v3)}
}

// Hash32 returns a 32-bit hash.
// The method decodes the value as four big-endian uint64s and then xors their halves.
func (b Bin256) Hash32() uint32 {
	v0 := binary.BigEndian.Uint64(b[:])
	v1 := binary.BigEndian.Uint64(b[8:])
	v2 := binary.BigEndian.Uint64(b[16:])
	v3 := binary.BigEndian.Uint64(b[24:])

	v0 = v0 ^ (v0 >> 32)
	v1 = v1 ^ (v1 >> 32)
	v2 = v2 ^ (v2 >> 32)
	v3 = v3 ^ (v3 >> 32)

	v := v0 ^ v1 ^ v2 ^ v3
	return uint32(v)
}

// Hash64 returns a 64-bit hash.
// The method decodes the value as four big-endian uint64s and then xors them.
func (b Bin256) Hash64() uint64 {
	v0 := binary.BigEndian.Uint64(b[:])
	v1 := binary.BigEndian.Uint64(b[8:])
	v2 := binary.BigEndian.Uint64(b[16:])
	v3 := binary.BigEndian.Uint64(b[24:])

	v := v0 ^ v1 ^ v2 ^ v3
	return v
}

// Size returns 32 bytes.
func (b Bin256) Size() int {
	return len(b)
}

// String returns a 64-char lower-case hex-encoded string.
func (b Bin256) String() string {
	buf := make([]byte, CharLen256)
	hex.Encode(buf, b[:])
	return string(buf)
}

// AppendHexTo appends a 64-char lower-case hex-encoded string to a buffer.
func (b Bin256) AppendHexTo(buf []byte) []byte {
	n := len(buf)
	n1 := n + CharLen256

	buf = buf[:n1]
	hex.Encode(buf[n:], b[:])
	return buf
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
	b0, err := Parse256(buf)
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

	b0, err := ParseString256(s)
	if err != nil {
		return err
	}

	*b = b0
	return nil
}
