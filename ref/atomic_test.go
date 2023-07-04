package ref

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testInterface interface {
	Set(int)
	Get() int
	Free()
}

type testObject struct{ v int }

func (o *testObject) Set(i int) { o.v = i }
func (o *testObject) Get() int  { return o.v }
func (o *testObject) Free()     { o.v = 0 }

func TestMap__should_map_referenced_object(t *testing.T) {
	obj := &testObject{10}
	r0 := New[*testObject](obj)
	r1 := Map[*testObject, testInterface](r0, func(o *testObject) testInterface {
		return o
	})

	assert.Equal(t, int64(2), r0.refs)
	assert.Equal(t, 10, r1.Unwrap().Get())

	r1.Release()
	assert.Equal(t, int64(1), r0.refs)
}
