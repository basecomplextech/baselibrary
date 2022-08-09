package bin128

import (
	"crypto/rand"
	"encoding/binary"
	"io"
	"time"
)

// Random returns a random B128.
func Random() B128 {
	return gen.random()
}

// TimeRandom returns a time-random B128.
func TimeRandom() B128 {
	return gen.timeRandom()
}

var gen = newGenerator()

type generator struct {
	rand io.Reader
}

func newGenerator() *generator {
	return &generator{
		rand: rand.Reader,
	}
}

func (g *generator) random() B128 {
	u := B128{}
	if _, err := g.rand.Read(u[:]); err != nil {
		panic(err)
	}
	return u
}

func (g *generator) timeRandom() B128 {
	u := B128{}

	now := g.timestamp()
	binary.BigEndian.PutUint64(u[:], now)

	if _, err := g.rand.Read(u[byteTimeLen:]); err != nil {
		panic(err)
	}
	return u
}

func (g *generator) timestamp() uint64 {
	now := time.Now()
	ts := now.UnixNano() / int64(time.Millisecond)
	return uint64(ts)
}
