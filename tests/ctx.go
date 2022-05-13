package tests

import (
	"reflect"
)

// X is a test case context.
type X struct {
	T

	values map[reflect.Type]interface{}
}

// NewX tries to cast T to a context or returns a new empty context.
func NewX(t T) *X {
	x, ok := t.(*X)
	if ok {
		return x
	}

	return &X{
		T: t,

		values: make(map[reflect.Type]interface{}),
	}
}

// Get returns a value from the context or inits a new one.
func Get[V any](x *X, init func(t T) V) (result V) {
	typ := reflect.TypeOf(result)

	v, ok := x.values[typ]
	if ok {
		return (v).(V)
	}

	val := init(x)
	x.values[typ] = val
	return val
}
