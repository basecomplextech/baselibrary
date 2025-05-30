// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package inject

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type Context interface {
	// Add adds multiple constructors or objects to the context, panics on duplicates.
	// When passed an object, the method adds the concrete object, not an interface.
	//
	// If you need to add an interface:
	//	- Wrap an object into inject.Typed(obj),
	//	- Or use generic inject.Add[T].
	Add(vv ...any) Context

	// AddOnce adds a constructor or an object if it does not exist already, skips duplicates.
	// When passed an object, the method adds the concrete object, not an interface.
	//
	// If you need to add an interface:
	//	- Wrap an object into inject.Typed(obj),
	//	- Or use generic inject.AddOnce[T].
	AddOnce(vv ...any) Context

	// AddObject adds an object (not a constructor) to the context, panics on duplicates.
	// The method treats functions as objects, not constructors.
	//
	// If you need to add an interface:
	//	- Wrap an object into inject.Typed(obj),
	//	- Or use generic inject.AddObject[T].
	AddObject(obj any) Context

	// Get inits an object and sets the pointer.
	Get(ptr any) Context
}

// New returns a new context with the provided constructors and objects.
func New(vv ...any) Context {
	x := newContext()
	x.Add(vv...)
	return x
}

// Add adds a constructor or an object to the context.
func Add[T any](c Context, v T) {
	typ := reflect.TypeFor[T]()
	p := newProvider(v)

	x := c.(*context)
	x.mu.Lock()
	defer x.mu.Unlock()

	x.addType(typ, p, false /* panic on duplicates */)
}

// AddOnce adds a constructor or an object if it does not exist already.
func AddOnce[T any](c Context, v T) {
	typ := reflect.TypeFor[T]()
	p := newProvider(v)

	x := c.(*context)
	x.mu.Lock()
	defer x.mu.Unlock()

	x.addType(typ, p, true /* skip existing */)
}

// AddObject adds an object to the context.
func AddObject[T any](c Context, obj T) {
	typ := reflect.TypeOf(obj)
	p := newObjectProvider(obj)

	x := c.(*context)
	x.mu.Lock()
	defer x.mu.Unlock()

	x.addType(typ, p, false /* panic on duplicates */)
}

// Get returns an object.
func Get[T any](c Context) T {
	var obj T
	c.Get(&obj)
	return obj
}

// internal

var _ Context = (*context)(nil)

type context struct {
	mu sync.Mutex

	objects   map[reflect.Type]reflect.Value
	providers map[reflect.Type]provider

	stack    []reflect.Type
	stackMap map[reflect.Type]struct{}
}

func newContext() *context {
	return &context{
		objects:   make(map[reflect.Type]reflect.Value),
		providers: make(map[reflect.Type]provider),
		stackMap:  make(map[reflect.Type]struct{}),
	}
}

// Add adds multiple constructors or objects to the context, panics on duplicates.
func (x *context) Add(vv ...any) Context {
	x.mu.Lock()
	defer x.mu.Unlock()

	for _, v := range vv {
		p := newProvider(v)
		x.add(p, false /* panic on duplicates */)
	}
	return x
}

// AddObject adds an object (not a constructor) to the context, panics on duplicates.
func (x *context) AddOnce(vv ...any) Context {
	x.mu.Lock()
	defer x.mu.Unlock()

	for _, v := range vv {
		p := newProvider(v)
		x.add(p, true /* skip existing */)
	}
	return x
}

// AddObject adds an object to the context.
func (x *context) AddObject(obj any) Context {
	p := newObjectProvider(obj)

	x.mu.Lock()
	defer x.mu.Unlock()

	x.add(p, false /* panic on duplicates */)
	return x
}

// Get inits an object and sets the pointer.
func (x *context) Get(ptr any) Context {
	pval := reflect.ValueOf(ptr)
	if pval.Kind() != reflect.Ptr {
		panic("must be a pointer")
	}

	elem := pval.Elem()
	typ := elem.Type()

	x.mu.Lock()
	defer x.mu.Unlock()

	obj := x.get(typ)
	elem.Set(obj)
	return x
}

// private

func (x *context) add(p provider, skipExisting bool) {
	typ := p.typ()

	// Check for duplicates
	_, ok := x.providers[typ]
	if ok {
		if skipExisting {
			return
		}
		panic(fmt.Sprintf("duplicate provider: %v", typ))
	}

	// Add provider
	x.providers[typ] = p
}

func (x *context) addType(typ reflect.Type, p provider, skipExisting bool) {
	_, ok := x.providers[typ]

	// Check for duplicates
	if ok {
		if skipExisting {
			return
		}
		panic(fmt.Sprintf("duplicate provider: %v", typ))
	}

	// Add provider
	x.providers[typ] = p
}

func (x *context) get(typ reflect.Type) reflect.Value {
	// Maybe return object
	obj, ok := x.objects[typ]
	if ok {
		return obj
	}

	// Add to stack, check for cycles
	_, cycle := x.stackMap[typ]
	x.stack = append(x.stack, typ)
	x.stackMap[typ] = struct{}{}

	// Panic if cycle
	if cycle {
		stack := x.stackString()
		panic(fmt.Sprintf("cycle detected: %v", stack))
	}

	// Get provider
	provider, ok := x.providers[typ]
	if !ok {
		stack := x.stackString()
		panic(fmt.Sprintf("provider not found: %v", stack))
	}

	// Init object
	obj = provider.init(x)
	x.objects[typ] = obj

	// Delete from stack
	x.stack = x.stack[:len(x.stack)-1]
	delete(x.stackMap, typ)
	return obj
}

// stackString returns the stack in direct order.
//
// Example: server.Server -> blobs.Blobs -> logging.Logger
func (x *context) stackString() string {
	sb := strings.Builder{}

	for i, s := range x.stack {
		if i != 0 {
			sb.WriteString(" -> ")
		}
		sb.WriteString(s.String())
	}

	return sb.String()
}
