package bin256

import (
	"crypto/rand"
	"encoding/binary"
	"io"
	"time"
)

// Random returns a random B256.
func Random() B256 {
	return gen.random()
}

// TimeRandom returns a time-random B128.
func TimeRandom() B256 {
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

func (g *generator) random() B256 {
	u := B256{}
	if _, err := g.rand.Read(u[:]); err != nil {
		panic(err)
	}
	return u
}

func (g *generator) timeRandom() B256 {
	u := B256{}

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
