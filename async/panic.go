package async

import (
	"fmt"
	"runtime/debug"
)

// Panic is an error which wraps a future panic.
// Panics must be passed as pointers.
type Panic struct {
	E     interface{}
	Stack []byte
}

func (p *Panic) Error() string {
	return fmt.Sprintf("panic: %v", p.E)
}

// Recover returns a panic if e is not nil, nil otherwise.
func Recover(e interface{}) *Panic {
	if e == nil {
		return nil
	}

	s := debug.Stack()
	return &Panic{E: e, Stack: s}
}
