package try

import (
	"fmt"
	"runtime/debug"
)

var _ error = (*Panic)(nil)

// Panic wraps a panic into an error, and includes a stack trace.
type Panic struct {
	E     interface{}
	Stack []byte
}

func (p *Panic) Error() string {
	return fmt.Sprintf("panic: %v", p.E)
}

// Recover returns a panic error if e not nil, includes the stack trace.
func Recover(e interface{}) *Panic {
	if e == nil {
		return nil
	}

	return RecoverStack(e, true)
}

// RecoverStack returns a panic error if e not nil, and optionally includes the stack trace.
func RecoverStack(e interface{}, stack bool) *Panic {
	if e == nil {
		return nil
	}

	var s []byte
	if stack {
		s = debug.Stack()
	}

	return &Panic{E: e, Stack: s}
}
