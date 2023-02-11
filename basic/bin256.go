// Package bin256 provides a binary 256-bit value.
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
	Bin256ByteLen = 32
	Bin256CharLen = (Bin256ByteLen * 2)
)

var Bin256Pattern = regexp.MustCompile(`^[0-9A-Za-z]{64}$`)

// Bin256 is a 32 byte value.
type Bin256 [Bin256ByteLen]byte

// Bin256FromInt converts an int into a big endian Bin256.
func Bin256FromInt(v int) Bin256 {
	b := Bin256{}
	buf := b[24:]

	binary.BigEndian.PutUint64(buf, uint64(v))
	return b
}

// Random

// RandomBin256 returns a random bin256.
func RandomBin256() Bin256 {
	u := Bin256{}
	if _, err := rand.Read(u[:]); err != nil {
		panic(err)
	}
	return u
}

// TimeRandomBin256 returns a time-random bin256 with a millisecond resolution.
func TimeRandomBin256() Bin256 {
	u := Bin256{}

	now := time.Now()
	ts := now.UnixNano() / int64(time.Millisecond)
	binary.BigEndian.PutUint64(u[:], uint64(ts))

	if _, err := rand.Read(u[8:]); err != nil {
		panic(err)
	}
	return u
}

// Parse

// ParseBin256 parses a bin256 from a 32-byte array.
func ParseBin256(b []byte) (Bin256, error) {
	switch {
	case b == nil:
		return Bin256{}, nil
	case len(b) == 0:
		return Bin256{}, nil
	case len(b) != Bin256ByteLen:
		return Bin256{}, errors.New("bin256: invalid bin256 length")
	}

	u := Bin256{}
	copy(u[:], b)
	return u, nil
}

// ParseBin256String parses a bin256 from 64-char string.
func ParseBin256String(s string) (Bin256, error) {
	switch {
	case s == "":
		return Bin256{}, nil
	case len(s) == 0:
		return Bin256{}, nil
	case len(s) != Bin256CharLen:
		return Bin256{}, errors.New("bin256: invalid bin256 length")
	}

	u := Bin256{}
	_, err := hex.Decode(u[:], []byte(s))
	if err != nil {
		return u, err
	}
	return u, nil
}

// MustParseBin256String parses a bin256 from 32-char string or panics.
func MustParseBin256String(s string) Bin256 {
	u, err := ParseBin256String(s)
	if err != nil {
		panic(err)
	}
	return u
}

// Sort

// SortBin256 sorts bin256 values.
func SortBin256(vv []Bin256) {
	sort.Slice(vv, func(i, j int) bool {
		a := vv[i]
		b := vv[j]
		return a.Less(b)
	})
}

// Methods

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

// Size returns 32 bytes.
func (b Bin256) Size() int {
	return len(b)
}

// String returns a 64-char lower-case hex-encoded string.
func (b Bin256) String() string {
	buf := make([]byte, Bin256CharLen)
	hex.Encode(buf, b[:])
	return string(buf)
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
	b0, err := ParseBin256(buf)
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

	b0, err := ParseBin256String(s)
	if err != nil {
		return err
	}

	*b = b0
	return nil
}
