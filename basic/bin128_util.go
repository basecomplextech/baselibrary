package basic

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

var Bin128Pattern = regexp.MustCompile(`^[0-9A-Za-z]{32}$`)

// Random

// RandomBin128 returns a random bin128.
func RandomBin128() Bin128 {
	return gen128.random()
}

// TimeRandomBin128 returns a time-random bin128 with a millisecond resolution.
func TimeRandomBin128() Bin128 {
	return gen128.timeRandom()
}

// Match

// MatchBin128 returns true if a byte string matches a bin128 pattern.
func MatchBin128(s []byte) bool {
	return Bin128Pattern.Match(s)
}

// MatchString128 returns true if a string matches a bin128 pattern.
func MatchString128(s string) bool {
	return Bin128Pattern.MatchString(s)
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

// ParseStringBin128 parses a bin128 from 33-char string.
func ParseStringBin128(s string) (Bin128, error) {
	return ParseByteStringBin128([]byte(s))
}

// ParseByteStringBin128 parses a bin128 from 32-char string.
func ParseByteStringBin128(s []byte) (Bin128, error) {
	switch {
	case s == nil:
		return Bin128{}, nil
	case len(s) == 0:
		return Bin128{}, nil
	case len(s) != Bin128CharLen:
		return Bin128{}, errors.New("bin128: invalid bin128 string length")
	}

	u := Bin128{}
	_, err := hex.Decode(u[:], s)
	if err != nil {
		return u, err
	}
	return u, nil
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

// private

var gen128 = newGenerator128()

type generator128 struct {
	rand io.Reader
}

func newGenerator128() *generator128 {
	return &generator128{
		rand: rand.Reader,
	}
}

func (g *generator128) random() Bin128 {
	u := Bin128{}
	if _, err := g.rand.Read(u[:]); err != nil {
		panic(err)
	}
	return u
}

func (g *generator128) timeRandom() Bin128 {
	u := Bin128{}

	now := g.timestamp()
	binary.BigEndian.PutUint64(u[:], now)

	if _, err := g.rand.Read(u[8:]); err != nil {
		panic(err)
	}
	return u
}

func (g *generator128) timestamp() uint64 {
	now := time.Now()
	ts := now.UnixNano() / int64(time.Millisecond)
	return uint64(ts)
}
