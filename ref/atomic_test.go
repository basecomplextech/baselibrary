package ref

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testInterface0 interface {
	Set(int)
	Get() int
	Free()
}

type testInterface1 interface {
	Get() int
	Free()
}

var _ testInterface0 = (*testValue)(nil)
var _ testInterface1 = (*testValue)(nil)

type testValue struct{ v int }

func (v *testValue) Set(i int) { v.v = i }
func (v *testValue) Get() int  { return v.v }
func (v *testValue) Free()     { v.v = 0 }

func TestCast__should_cast_interface_and_share_refcount(t *testing.T) {
	v := &testValue{10}
	r0 := Wrap[testInterface0](v)
	r1 := Cast[testInterface0, testInterface1](r0)

	assert.Same(t, r0.obj, r1.obj)
	assert.Same(t, r0.refs, r1.refs)
	assert.Equal(t, int64(0), r1._refs)
	assert.Equal(t, 10, r1.Unwrap().Get())

	r0.Release()
	assert.Equal(t, int64(0), r0.Refcount())
	assert.Equal(t, int64(0), r1.Refcount())
	assert.Equal(t, 0, v.v)
}
