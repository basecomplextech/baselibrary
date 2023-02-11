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
	Bin64ByteLen = 8
	Bin64CharLen = (Bin64ByteLen * 2) // 341a7d60bc5893a6
)

var Bin64Pattern = regexp.MustCompile(`^[0-9A-Za-z]{16}$`)

// Bin64 is a binary 64-bit value.
type Bin64 [Bin64ByteLen]byte

// Bin64FromInt converts an int into a bin64.
func Bin64FromInt(v int) Bin64 {
	b := Bin64{}
	binary.BigEndian.PutUint64(b[:], uint64(v))
	return b
}

// Random

// RandomBin64 returns a random bin64.
func RandomBin64() Bin64 {
	u := Bin64{}
	if _, err := rand.Read(u[:]); err != nil {
		panic(err)
	}
	return u
}

// TimeRandomBin64 returns a time-random bin64 with a second resolution.
func TimeRandomBin64() Bin64 {
	u := Bin64{}

	now := time.Now()
	ts := now.UnixNano() / int64(time.Second)
	binary.BigEndian.PutUint32(u[:], uint32(ts))

	if _, err := rand.Read(u[4:]); err != nil {
		panic(err)
	}
	return u
}

// Parse

// ParseBin64 parses a bin64 value from a 16-byte array.
func ParseBin64(b []byte) (Bin64, error) {
	switch {
	case b == nil:
		return Bin64{}, nil
	case len(b) == 0:
		return Bin64{}, nil
	case len(b) != Bin64ByteLen:
		return Bin64{}, errors.New("bin64: invalid bin64 length")
	}

	u := Bin64{}
	copy(u[:], b)
	return u, nil
}

// ParseBin64String parses a bin64 from 32-char string.
func ParseBin64String(s string) (Bin64, error) {
	switch {
	case s == "":
		return Bin64{}, nil
	case len(s) == 0:
		return Bin64{}, nil
	case len(s) != Bin64CharLen:
		return Bin64{}, errors.New("bin64: invalid bin64 string length")
	}

	u := Bin64{}
	_, err := hex.Decode(u[:], []byte(s))
	if err != nil {
		return u, err
	}
	return u, nil
}

// Sort

// SortBin64 sorts bin64 values.
func SortBin64(vv []Bin64) {
	sort.Slice(vv, func(i, j int) bool {
		a := vv[i]
		b := vv[j]
		return a.Less(b)
	})
}

// Methods

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
	buf := make([]byte, Bin64CharLen)
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
	b0, err := ParseBin64(buf)
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

	b0, err := ParseBin64String(s)
	if err != nil {
		return err
	}

	*b = b0
	return nil
}
