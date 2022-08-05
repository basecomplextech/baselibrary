package tests

import (
	"fmt"
	"reflect"
)

// X is a test case context.
type X struct {
	T

	objects map[reflect.Type]any
	funcs   map[reflect.Type]func() any
}

// NewContext returns a new test context.
func NewContext(t T) *X {
	return &X{
		T: t,

		funcs:   make(map[reflect.Type]func() any),
		objects: make(map[reflect.Type]any),
	}
}

// Add adds a constructor to the context or panics on duplicate constructor.
func Add[T any](x *X, fn func() T) {
	var ptr *T // ptr for interfaces
	typ := reflect.TypeOf(ptr).Elem()

	x.add(typ, func() any {
		return fn()
	})
}

// Get returns an object from the context or panics.
func Get[T any](x *X) T {
	var ptr *T // ptr for interfaces
	typ := reflect.TypeOf(ptr).Elem()

	o, ok := x.objects[typ]
	if ok {
		return (o).(T)
	}

	p, ok := x.funcs[typ]
	if !ok {
		panic(fmt.Sprintf("object/constructor not found %v", typ))
	}

	o = p()
	return (o).(T)
}

// Add adds a constructor to the context or panics on duplicate constructor.
// The constructor signature must be `func() T`.
func (x *X) Add(fn any) {
	f := reflect.TypeOf(fn)
	if f.Kind() != reflect.Func {
		panic(fmt.Sprintf("constructor must be function, %v", f))
	}
	if f.NumIn() != 0 {
		panic(fmt.Sprintf("constructor must be empty function, %v", f))
	}
	if f.NumOut() != 1 {
		panic(fmt.Sprintf("constructor must return one value, %v", f))
	}

	typ := f.Out(0)
	x.add(typ, func() any {
		out := reflect.ValueOf(fn).Call(nil)[0]
		return out.Interface()
	})
}

// private

func (x *X) add(typ reflect.Type, fn func() any) {
	_, ok := x.objects[typ]
	if ok {
		panic(fmt.Sprintf("object already initialized %v", typ))
	}

	_, ok = x.funcs[typ]
	if ok {
		panic(fmt.Sprintf("duplicate object constructor %v", typ))
	}

	x.funcs[typ] = fn
}
