package bin

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

var Pattern128 = regexp.MustCompile(`^[0-9A-Za-z]{32}$`)

// Random

// Random128 returns a random bin128.
func Random128() Bin128 {
	return gen128.random()
}

// TimeRandom128 returns a time-random bin128.
func TimeRandom128() Bin128 {
	return gen128.timeRandom()
}

// Match

// Match128 returns true if a byte string matches a bin128 pattern.
func Match128(s []byte) bool {
	return Pattern128.Match(s)
}

// MatchString128 returns true if a string matches a bin128 pattern.
func MatchString128(s string) bool {
	return Pattern128.MatchString(s)
}

// Parse

// Parse128 parses a bin128 value from a 16-byte array.
func Parse128(b []byte) (Bin128, error) {
	switch {
	case b == nil:
		return Bin128{}, nil
	case len(b) == 0:
		return Bin128{}, nil
	case len(b) != ByteLen128:
		return Bin128{}, errors.New("bin128: invalid bin128 length")
	}

	u := Bin128{}
	copy(u[:], b)
	return u, nil
}

// ParseString128 parses a bin128 from 33-char string.
func ParseString128(s string) (Bin128, error) {
	return ParseByteString128([]byte(s))
}

// ParseByteString128 parses a bin128 from 32-char string.
func ParseByteString128(s []byte) (Bin128, error) {
	switch {
	case s == nil:
		return Bin128{}, nil
	case len(s) == 0:
		return Bin128{}, nil
	case len(s) != CharLen128:
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
