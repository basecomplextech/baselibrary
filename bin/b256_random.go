package bin

import (
	"crypto/rand"
	"encoding/binary"
	"io"
	"time"
)

// Random256 returns a random bin256.
func Random256() Bin256 {
	return gen256.random()
}

// TimeRandom256 returns a time-random bin256.
func TimeRandom256() Bin256 {
	return gen256.timeRandom()
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
