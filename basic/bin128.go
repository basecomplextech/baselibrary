package basic

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"regexp"
	"sort"
	"time"
)

const (
	Bin128ByteLen = 16
	Bin128CharLen = (Bin128ByteLen * 2) // 341a7d60bc5893a64bda3de06721534c
)

var Bin128Pattern = regexp.MustCompile(`^[0-9A-Za-z]{32}$`)

// Bin128 is a 128-bit value.
type Bin128 [Bin128ByteLen]byte

// Bin128FromInt converts an int into a bin128.
func Bin128FromInt(v int) Bin128 {
	b := Bin128{}
	buf := b[8:]

	binary.BigEndian.PutUint64(buf, uint64(v))
	return b
}

// Random

// RandomBin128 returns a random bin128.
func RandomBin128() Bin128 {
	u := Bin128{}
	if _, err := rand.Read(u[:]); err != nil {
		panic(err)
	}
	return u
}

// TimeRandomBin128 returns a time-random bin128 with a millisecond resolution.
func TimeRandomBin128() Bin128 {
	u := Bin128{}

	now := time.Now()
	ts := now.UnixNano() / int64(time.Millisecond)
	binary.BigEndian.PutUint64(u[:], uint64(ts))

	if _, err := rand.Read(u[8:]); err != nil {
		panic(err)
	}
	return u
}

// Sort

// SortBin128 sorts bin128 values.
func SortBin128(vv []Bin128) {
	sort.Slice(vv, func(i, j int) bool {
		a := vv[i]
		b := vv[j]
		return a.Less(b)
	})
}

// Parse

// ParseBin128 parses a bin128 value from a 16-byte array.
func ParseBin128(b []byte) (Bin128, error) {
	switch {
	case b == nil:
		return Bin128{}, nil
	case len(b) == 0:
		return Bin128{}, nil
	case len(b) != Bin128ByteLen:
		return Bin128{}, errors.New("bin128: invalid bin128 length")
	}

	u := Bin128{}
	copy(u[:], b)
	return u, nil
}

// ParseBin128String parses a bin128 from 32-char string.
func ParseBin128String(s string) (Bin128, error) {
	switch {
	case s == "":
		return Bin128{}, nil
	case len(s) == 0:
		return Bin128{}, nil
	case len(s) != Bin128CharLen:
		return Bin128{}, errors.New("bin128: invalid bin128 string length")
	}

	u := Bin128{}
	_, err := hex.Decode(u[:], []byte(s))
	if err != nil {
		return u, err
	}
	return u, nil
}

// Methods

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
	buf := make([]byte, Bin128CharLen)
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
	b0, err := ParseBin128(buf)
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

	b0, err := ParseBin128String(s)
	if err != nil {
		return err
	}

	*b = b0
	return nil
}
