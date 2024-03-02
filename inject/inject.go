package inject

import (
	"reflect"
	"sync"
)

type Context struct {
	mu  sync.Mutex
	ctx *context
}

// New returns a new context.
func New() *Context {
	return &Context{
		ctx: newContext(),
	}
}

// Add adds a constructor or an object to the context.
func Add[T any](c *Context, v T) {
	typ := reflect.TypeOf(v)
	p := newProvider(v)

	c.mu.Lock()
	defer c.mu.Unlock()

	c.ctx.addType(typ, p, false)
}

// AddOnce adds a constructor or an object if it does not exist already.
func AddOnce[T any](c *Context, v T) {
	var zero *T
	typ := reflect.TypeOf(zero).Elem()
	p := newProvider(v)

	c.mu.Lock()
	defer c.mu.Unlock()

	c.ctx.addType(typ, p, true /* skip existing */)
}

// AddObject adds an object to the context.
func AddObject[T any](c *Context, obj T) {
	typ := reflect.TypeOf(obj)
	p := newObjectProvider(obj)

	c.mu.Lock()
	defer c.mu.Unlock()

	c.ctx.addType(typ, p, false)
}

// Get returns an object.
func Get[T any](c *Context) T {
	var obj T
	c.Get(&obj)
	return obj
}

// Add adds a constructor or an object to the context.
//
// When passed an object, the method adds the concrete object, not an interface.
// Should you need to add an interface, use inject.Add[T], or pass a function
// that returns an interface `func() T { return obj }`.
func (c *Context) Add(v any) *Context {
	p := newProvider(v)

	c.mu.Lock()
	defer c.mu.Unlock()

	c.ctx.add(p, false)
	return c
}

// AddOnce adds a constructor or an object if it does not exist already.
//
// When passed an object, the method adds the concrete object, not an interface.
// Should you need to add an interface, use inject.AddOnce[T], or pass a function
// that returns an interface `func() T { return obj }`.
func (c *Context) AddOnce(v any) *Context {
	p := newProvider(v)

	c.mu.Lock()
	defer c.mu.Unlock()

	c.ctx.add(p, true /* skip existing */)
	return c
}

// AddObject adds an object to the context.
//
// When passed an object, the method adds the concrete object, not an interface.
// Should you need to add an interface, use inject.AddObject[T].
func (c *Context) AddObject(obj any) *Context {
	p := newObjectProvider(obj)

	c.mu.Lock()
	defer c.mu.Unlock()

	c.ctx.add(p, false)
	return c
}

// Get inits an object and sets the pointer.
func (c *Context) Get(ptr any) *Context {
	pval := reflect.ValueOf(ptr)
	if pval.Kind() != reflect.Ptr {
		panic("must be a pointer")
	}

	elem := pval.Elem()
	typ := elem.Type()

	c.mu.Lock()
	defer c.mu.Unlock()

	obj := c.ctx.get(typ)
	elem.Set(obj)
	return c
}
