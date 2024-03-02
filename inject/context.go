package inject

import (
	"fmt"
	"reflect"
	"strings"
)

type context struct {
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
