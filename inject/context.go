package inject

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type Context interface {
	// Add adds a constructor or an object to the context.
	//
	// When passed an object, the method adds the concrete object, not an interface.
	// Should you need to add an interface, use inject.Add[T], or pass a function
	// that returns an interface `func() T { return obj }`.
	Add(v any) Context

	// AddOnce adds a constructor or an object if it does not exist already.
	//
	// When passed an object, the method adds the concrete object, not an interface.
	// Should you need to add an interface, use inject.AddOnce[T], or pass a function
	// that returns an interface `func() T { return obj }`.
	AddOnce(v any) Context

	// AddObject adds an object to the context.
	//
	// When passed an object, the method adds the concrete object, not an interface.
	// Should you need to add an interface, use inject.AddObject[T].
	AddObject(obj any) Context

	// Get inits an object and sets the pointer.
	Get(ptr any) Context
}

// New returns a new context.
func New() Context {
	return newContext()
}

// Add adds a constructor or an object to the context.
func Add[T any](c Context, v T) {
	typ := reflect.TypeFor[T]()
	p := newProvider(v)

	x := c.(*context)
	x.mu.Lock()
	defer x.mu.Unlock()

	x.addType(typ, p, false)
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

	x.addType(typ, p, false)
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

// Add adds a constructor or an object to the context.
func (x *context) Add(v any) Context {
	p := newProvider(v)

	x.mu.Lock()
	defer x.mu.Unlock()

	x.add(p, false)
	return x
}

// AddOnce adds a constructor or an object if it does not exist already.
func (x *context) AddOnce(v any) Context {
	p := newProvider(v)

	x.mu.Lock()
	defer x.mu.Unlock()

	x.add(p, true /* skip existing */)
	return x
}

// AddObject adds an object to the context.
func (x *context) AddObject(obj any) Context {
	p := newObjectProvider(obj)

	x.mu.Lock()
	defer x.mu.Unlock()

	x.add(p, false)
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

// stackString returns the stack in reverse order.
//
// Example: logging.Logger <- blobs.Blobs <- server.Server
func (x *context) stackString() string {
	sb := strings.Builder{}

	for i := len(x.stack) - 1; i >= 0; i-- {
		sb.WriteString(x.stack[i].String())

		if i > 0 {
			sb.WriteString(" <- ")
		}
	}

	return sb.String()
}
