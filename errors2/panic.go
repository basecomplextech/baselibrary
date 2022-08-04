package errors2

import (
	"fmt"
	"runtime/debug"
)

// PanicError wraps a panic into an error, and optionally includes a stack trace.
type PanicError struct {
	E     interface{}
	Stack []byte
}

// Error returns the error message.
func (e *PanicError) Error() string {
	return fmt.Sprintf("%v", e.E)
}

// Recover wraps a panic into a *PanicError if e is not nil, includes the stack trace.
func Recover(e interface{}) error {
	if e == nil {
		return nil
	}

	return &PanicError{E: e}
}

// RecoverStack wraps a panic into a *PanicError if e is not nil, and includes the stack trace.
func RecoverStack(e interface{}) (err error, stack []byte) {
	if e == nil {
		return nil, nil
	}

	stack = debug.Stack()
	err = &PanicError{E: e, Stack: stack}
	return
}
