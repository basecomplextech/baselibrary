package u128

import (
	"crypto/rand"
	"encoding/binary"
	"io"
	"time"
)

// Generator generates random and time-random U128s.
type Generator interface {
	// RandomU128 returns a random U128.
	RandomU128() U128

	// TimeRandom128 returns a time-random U128.
	TimeRandom128() U128
}

var global Generator = newGenerator()

type generator struct {
	rand io.Reader
}

func newGenerator() *generator {
	return &generator{
		rand: rand.Reader,
	}
}

// RandomU128 returns a random U128.
func (g *generator) RandomU128() U128 {
	u := U128{}

	if _, err := g.rand.Read(u[:]); err != nil {
		panic(err)
	}

	return u
}

// TimeRandom128 returns a time-random U128.
func (g *generator) TimeRandom128() U128 {
	u := U128{}

	ts := g.timestamp()
	binary.BigEndian.PutUint64(u[:], ts)

	if _, err := g.rand.Read(u[byteTimeLen:]); err != nil {
		panic(err)
	}

	return u
}

func (i *generator) timestamp() uint64 {
	now := time.Now()
	ts := now.UnixNano() / int64(time.Millisecond)
	return uint64(ts)
}
