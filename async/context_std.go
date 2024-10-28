// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import (
	context_ "context"
	"time"

	"github.com/basecomplextech/baselibrary/status"
)

// ContextToStd returns a standard library context from an async one.
func ContextToStd(ctx Context) context_.Context {
	return newStdContext(ctx)
}

// internal

var _ context_.Context = (*stdContext)(nil)

type stdContext struct {
	ctx Context
}

func newStdContext(ctx Context) *stdContext {
	return &stdContext{ctx: ctx}
}

// Deadline returns the time when work done on behalf of this context should be canceled.
func (x *stdContext) Deadline() (deadline time.Time, ok bool) {
	return
}

// Done returns a channel that's closed when work done on behalf of this context should be canceled.
func (x *stdContext) Done() <-chan struct{} {
	return x.ctx.Wait()
}

// If Done is not yet closed, Err returns nil.
func (x *stdContext) Err() error {
	st := x.ctx.Status()
	switch st.Code {
	case status.CodeCancelled:
		return context_.Canceled
	case status.CodeTimeout:
		return context_.DeadlineExceeded
	}
	return status.ToError(st)
}

// Value returns the value associated with this context for key, or nil
// if no value is associated with key.
func (x *stdContext) Value(key any) any {
	return nil
}
