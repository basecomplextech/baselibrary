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

	return RecoverStack(e, true)
}

// RecoverStack wraps a panic into a *PanicError if e is not nil, and optionally includes the stack trace.
func RecoverStack(e interface{}, stack bool) error {
	if e == nil {
		return nil
	}

	var s []byte
	if stack {
		s = debug.Stack()
	}

	return &PanicError{E: e, Stack: s}
}
