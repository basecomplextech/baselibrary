package di

import (
	"fmt"
	"reflect"
)

// X is a context which holds providers and initialized objects.
type X struct {
	objects   map[reflect.Type]any
	providers map[reflect.Type]func() any
}

// Add adds a provider to a context or panics on a duplicate provider.
func Add[T any](x *X, provider func() T) {
	var ptr *T // pointer to get nil interface type
	typ := reflect.TypeOf(ptr).Elem()

	x.add(typ, func() any {
		return provider()
	})
}

// Get returns an object from a context or panics on an absent provider.
func Get[T any](x *X) T {
	var ptr *T // pointer to get nil interface type
	typ := reflect.TypeOf(ptr).Elem()

	v := x.get(typ)
	return v.(T)
}

// Build builds an object from modules.
func Build[T any](modules ...func(x *X)) T {
	x := buildContext(modules...)
	return Get[T](x)
}

// BuildContext builds a context from modules.
func BuildContext[T any](modules ...func(*X)) *X {
	return buildContext(modules...)
}

// internal

func newContext() *X {
	return &X{
		objects:   make(map[reflect.Type]any),
		providers: make(map[reflect.Type]func() any),
	}
}

func buildContext(modules ...func(*X)) *X {
	x := newContext()
	for _, module := range modules {
		module(x)
	}
	return x
}

func (x *X) add(typ reflect.Type, provider func() any) {
	_, ok := x.objects[typ]
	if ok {
		panic(fmt.Sprintf("object already initialized %v", typ))
	}

	_, ok = x.providers[typ]
	if ok {
		panic(fmt.Sprintf("duplicate provider %v", typ))
	}

	x.providers[typ] = provider
}

func (x *X) get(typ reflect.Type) any {
	obj, ok := x.objects[typ]
	if ok {
		return obj
	}

	p, ok := x.providers[typ]
	if !ok {
		panic(fmt.Sprintf("provider not found %v", typ))
	}

	obj = p()
	x.objects[typ] = obj
	return obj
}
