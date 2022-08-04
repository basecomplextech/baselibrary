package di

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type A struct {
	b *B
}

type B struct {
	c *C
}

type C struct {
	d D
}

type D interface {
	String() string
}

type DStruct struct {
	s string
}

func (d *DStruct) String() string {
	return d.s
}

func newA(x *X) *A {
	return &A{
		b: Get[*B](x),
	}
}

func newB(x *X) *B {
	return &B{
		c: Get[*C](x),
	}
}

func newC(x *X) *C {
	return &C{
		d: Get[D](x),
	}
}

func TestBuild(t *testing.T) {
	a := Build[*A](func(x *X) {
		Add(x, func() *A {
			return newA(x)
		})

		Add(x, func() *B {
			return newB(x)
		})

		Add(x, func() *C {
			return newC(x)
		})

		Add(x, func() D {
			return &DStruct{"hello, world"}
		})
	})

	s := a.b.c.d.String()
	assert.Equal(t, "hello, world", s)
}
