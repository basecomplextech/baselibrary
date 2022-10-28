package types

import (
	"crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"regexp"
	"sort"
	"time"
)

var Bin256Pattern = regexp.MustCompile(`^[0-9A-Za-z]{64}$`)

// Random

// RandomBin256 returns a random bin256.
func RandomBin256() Bin256 {
	return gen256.random()
}

// TimeRandomBin256 returns a time-random bin256 with a millisecond resolution.
func TimeRandomBin256() Bin256 {
	return gen256.timeRandom()
}

// Match

// MatchBin256 returns true if a byte string matches a bin256 pattern.
func MatchBin256(s []byte) bool {
	return Bin256Pattern.Match(s)
}

// MatchStringBin256 returns true if a string matches a bin256 pattern.
func MatchStringBin256(s string) bool {
	return Bin256Pattern.MatchString(s)
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

// ParseStringBin256 parses a bin256 from 64-char string.
func ParseStringBin256(s string) (Bin256, error) {
	return ParseByteStringBin256([]byte(s))
}

// ParseByteStringBin256 parses a bin256 from 64-char byte string.
func ParseByteStringBin256(s []byte) (Bin256, error) {
	switch {
	case s == nil:
		return Bin256{}, nil
	case len(s) == 0:
		return Bin256{}, nil
	case len(s) != Bin256CharLen:
		return Bin256{}, errors.New("bin256: invalid bin256 length")
	}

	u := Bin256{}
	_, err := hex.Decode(u[:], s)
	if err != nil {
		return u, err
	}
	return u, nil
}

// MustParseStringBin256 parses a bin256 from 32-char string or panics.
func MustParseStringBin256(s string) Bin256 {
	u, err := ParseStringBin256(s)
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

// private

var gen256 = newGenerator256()

type generator256 struct {
	rand io.Reader
}

func newGenerator256() *generator256 {
	return &generator256{
		rand: rand.Reader,
	}
}

func (g *generator256) random() Bin256 {
	u := Bin256{}
	if _, err := g.rand.Read(u[:]); err != nil {
		panic(err)
	}
	return u
}

func (g *generator256) timeRandom() Bin256 {
	u := Bin256{}

	now := g.timestamp()
	binary.BigEndian.PutUint64(u[:], now)

	if _, err := g.rand.Read(u[8:]); err != nil {
		panic(err)
	}
	return u
}

func (g *generator256) timestamp() uint64 {
	now := time.Now()
	ts := now.UnixNano() / int64(time.Millisecond)
	return uint64(ts)
}
