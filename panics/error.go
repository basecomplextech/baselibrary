package panics

import (
	"fmt"
	"runtime/debug"
)

// Error wraps a panic into an error, by default includes the stack trace.
type Error struct {
	E     interface{}
	Stack []byte
}

// Error returns the error message.
func (e *Error) Error() string {
	return fmt.Sprintf("%v", e.E)
}

// Recover wraps a panic into a *Error if e is not nil, includes the stack trace.
func Recover(e interface{}) error {
	if e == nil {
		return nil
	}

	stack := debug.Stack()
	return &Error{E: e, Stack: stack}
}

// RecoverStack wraps a panic into a *Error if e is not nil, includes the stack trace.
func RecoverStack(e interface{}) (err error, stack []byte) {
	if e == nil {
		return nil, nil
	}

	stack = debug.Stack()
	err = &Error{E: e, Stack: stack}
	return
}

// RecoverNoStack wraps a panic into a *Error if e is not nil, skips the stack trace.
func RecoverNoStack(e interface{}) error {
	if e == nil {
		return nil
	}

	return &Error{E: e}
}
