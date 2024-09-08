// Copyright 2023 Ivan Korobkov. All rights reserved.

package panics

import (
	"fmt"
	"runtime/debug"
)

// Error wraps a panic into an error, by default includes the stack trace.
type Error struct {
	E     any
	Stack []byte
}

// Error returns the error message.
func (e *Error) Error() string {
	return fmt.Sprintf("%v", e.E)
}

// Panicf panics with a formatted message.
func Panicf(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	panic(msg)
}

// Recover wraps a panic into a *Error if e is not nil, includes the stack trace.
func Recover(e any) error {
	if e == nil {
		return nil
	}

	stack := debug.Stack()
	return &Error{E: e, Stack: stack}
}

// RecoverStack wraps a panic into a *Error if e is not nil, includes the stack trace.
func RecoverStack(e any) (err error, stack []byte) {
	if e == nil {
		return nil, nil
	}

	stack = debug.Stack()
	err = &Error{E: e, Stack: stack}
	return
}

// RecoverNoStack wraps a panic into a *Error if e is not nil, skips the stack trace.
func RecoverNoStack(e any) error {
	if e == nil {
		return nil
	}

	return &Error{E: e}
}
