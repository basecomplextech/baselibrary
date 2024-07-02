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
	CharLen128 = (ByteLen128 * 2) + 1 // 341a7d60bc5893a6-4bda3de06721534c
)

var (
	Max128 = Bin128{
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
	}
	Pattern128 = regexp.MustCompile(`^[0-9A-Za-z]{16}-[0-9A-Za-z]{16}$`)
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

// Join128 joins two bin64 values into a bin128.
func Join128(b0, b1 Bin64) Bin128 {
	b := Bin128{}
	copy(b[:], b0[:])
	copy(b[8:], b1[:])
	return b
}

// Compare compares two values.
//
// Returns:
//
//	-1 if a < b
//	 0 if a == b
//	 1 if a > b
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

// Parts

// Ints returns two int64s decoded as big-endian.
func (b Bin128) Ints() (int64, int64) {
	v0 := binary.BigEndian.Uint64(b[:])
	v1 := binary.BigEndian.Uint64(b[8:])
	return int64(v0), int64(v1)
}

// Parts returns two bin64 values.
func (b Bin128) Parts() [2]Bin64 {
	var b0, b1 Bin64
	copy(b0[:], b[:])
	copy(b1[:], b[8:])
	return [2]Bin64{b0, b1}
}

// Hash

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

// String/Hex

// String returns a 33-char lower-case hex-encoded string.
func (b Bin128) String() string {
	buf := make([]byte, CharLen128)
	hex.Encode(buf, b[:8])
	buf[16] = '-'
	hex.Encode(buf[17:], b[8:])
	return string(buf)
}

// AppendHexTo appends a 33-char lower-case hex-encoded string to a buffer.
func (b Bin128) AppendHexTo(buf []byte) []byte {
	n := len(buf)
	n1 := n + CharLen128

	buf = buf[:n1]
	hex.Encode(buf[n:], b[:8])
	buf[n+16] = '-'
	hex.Encode(buf[n+17:], b[8:])
	return buf
}

// Marshal

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
