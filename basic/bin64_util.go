package basic

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"regexp"
	"sort"
	"time"
)

var Bin64Pattern = regexp.MustCompile(`^[0-9A-Za-z]{16}$`)

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

// Match

// MatchBin64 returns true if a byte string matches a bin64 pattern.
func MatchBin64(s []byte) bool {
	return Bin64Pattern.Match(s)
}

// MatchString64 returns true if a string matches a bin64 pattern.
func MatchString64(s string) bool {
	return Bin64Pattern.MatchString(s)
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

// ParseStringBin64 parses a bin64 from 33-char string.
func ParseStringBin64(s string) (Bin64, error) {
	return ParseByteStringBin64([]byte(s))
}

// ParseByteStringBin64 parses a bin64 from 32-char string.
func ParseByteStringBin64(s []byte) (Bin64, error) {
	switch {
	case s == nil:
		return Bin64{}, nil
	case len(s) == 0:
		return Bin64{}, nil
	case len(s) != Bin64CharLen:
		return Bin64{}, errors.New("bin64: invalid bin64 string length")
	}

	u := Bin64{}
	_, err := hex.Decode(u[:], s)
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
