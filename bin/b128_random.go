package bin

import (
	"crypto/rand"
	"encoding/binary"
	"io"
	"time"
)

// Random128 returns a random bin128.
func Random128() Bin128 {
	return gen128.random()
}

// TimeRandom128 returns a time-random bin128.
func TimeRandom128() Bin128 {
	return gen128.timeRandom()
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
