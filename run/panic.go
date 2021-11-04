package run

import (
	"fmt"
	"runtime/debug"
)

var _ error = (*Panic)(nil)

// Panic wraps a panic into an error.
type Panic struct {
	E     interface{}
	Stack []byte
}

func (p *Panic) Error() string {
	return fmt.Sprintf("panic: %v", p.E)
}

// Recover returns a panic error if e not nil.
func Recover(e interface{}) *Panic {
	if e == nil {
		return nil
	}

	return &Panic{E: e}
}

// RecoverStack returns a panic error if e not nil and includes the stack.
func RecoverStack(e interface{}) *Panic {
	if e == nil {
		return nil
	}

	s := debug.Stack()
	return &Panic{E: e, Stack: s}
}

// TryRecover calls a function, recovers and returns its panic as an error if any.
func TryRecover(fn func() error) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = Recover(e)
		}
	}()

	return fn()
}
