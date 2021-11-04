package u256

import (
	"crypto/rand"
	"io"
)

// Generator generates random and time-random U256s.
type Generator interface {
	// RandomU256 returns a random U256.
	RandomU256() U256
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

// RandomU256 returns a random U256.
func (g *generator) RandomU256() U256 {
	u := U256{}

	if _, err := g.rand.Read(u[:]); err != nil {
		panic(err)
	}

	return u
}
