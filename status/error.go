// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package status

import (
	"fmt"

	"github.com/basecomplextech/baselibrary/panics"
)

// Error

// Error returns an internal error status.
func Error(msg string) Status {
	return Status{
		Code:    CodeError,
		Message: msg,
	}
}

// Errorf formats and returns an internal error status.
func Errorf(format string, a ...any) Status {
	msg := format
	if len(a) > 0 {
		msg = fmt.Sprintf(format, a...)
	}

	return Status{
		Code:    CodeError,
		Message: msg,
	}
}

// WrapError returns an internal error status.
func WrapError(err error) Status {
	return wrapError(err, CodeError)
}

// WrapErrorf formats and returns an internal error status.
func WrapErrorf(err error, format string, a ...any) Status {
	return wrapErrorf(err, CodeError, format, a...)
}

// ExternalError

// ExternalError returns an external error status.
func ExternalError(msg string) Status {
	return Status{
		Code:    CodeExternalError,
		Message: msg,
	}
}

// ExternalErrorf formats and returns an external error status.
func ExternalErrorf(format string, a ...any) Status {
	msg := format
	if len(a) > 0 {
		msg = fmt.Sprintf(format, a...)
	}

	return Status{
		Code:    CodeExternalError,
		Message: msg,
	}
}

// WrapExternalError returns an external error status.
func WrapExternalError(err error) Status {
	msg := "External error"
	if err != nil {
		msg = err.Error()
	}

	return Status{
		Code:    CodeExternalError,
		Message: msg,
		Error:   err,
	}
}

// WrapExternalErrorf formats and returns an external error status.
func WrapExternalErrorf(err error, format string, a ...any) Status {
	msg := format
	if len(a) > 0 {
		msg = fmt.Sprintf(format, a...)
	}
	if err != nil {
		msg += ": " + err.Error()
	}

	return Status{
		Code:    CodeExternalError,
		Message: msg,
		Error:   err,
	}
}

// Recover

// Recover recovers from a panic and returns an internal error status.
func Recover(e any) Status {
	err := panics.Recover(e)
	return WrapError(err)
}

// RecoverStack recovers from a panic and returns an internal error status and a stack trace.
func RecoverStack(e any) (Status, []byte) {
	err, stack := panics.RecoverStack(e)
	return WrapError(err), stack
}

// internal

func wrapError(err error, code Code) Status {
	switch e := err.(type) {
	case nil:
		return OK
	case *Err:
		return e.Status()
	}

	return Status{
		Code:    code,
		Message: err.Error(),
		Error:   err,
	}
}

func wrapErrorf(err error, code Code, format string, a ...any) Status {
	msg := format
	if len(a) > 0 {
		msg = fmt.Sprintf(format, a...)
	}
	if err != nil {
		msg += ": " + err.Error()
	}

	e, ok := err.(*Err)
	if ok {
		return Status{
			Code:    e.Code,
			Message: msg,
			Error:   e.Cause,
		}
	}

	return Status{
		Code:    code,
		Message: msg,
		Error:   err,
	}
}
