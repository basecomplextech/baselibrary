// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import (
	context_ "context"
	"time"

	"github.com/basecomplextech/baselibrary/async/internal/context"
)

type (
	// Context is an async cancellation context.
	Context = context.Context

	// CancelContext is a cancellable async context.
	CancelContext = context.CancelContext

	// ContextCallback is called when the context is cancelled.
	ContextCallback = context.Callback
)

// New

// NewContext returns a new cancellable context.
func NewContext() CancelContext {
	return context.New()
}

// NoContext returns a non-cancellable background context.
func NoContext() Context {
	return context.No()
}

// CancelledContext returns a cancelled context.
func CancelledContext() CancelContext {
	return context.Cancelled()
}

// Timeout

// TimeoutContext returns a context with a timeout.
func TimeoutContext(timeout time.Duration) Context {
	return context.Timeout(timeout)
}

// DeadlineContext returns a context with a deadline.
func DeadlineContext(deadline time.Time) Context {
	return context.Deadline(deadline)
}

// Next

// NextContext returns a child context.
func NextContext(parent Context) CancelContext {
	return context.Next(parent)
}

// NextTimeoutContext returns a child context with a timeout.
func NextTimeoutContext(parent Context, timeout time.Duration) Context {
	return context.NextTimeout(parent, timeout)
}

// NextDeadlineContext returns a child context with a deadline.
func NextDeadlineContext(parent Context, deadline time.Time) Context {
	return context.NextDeadline(parent, deadline)
}

// Standard

// StdContext returns a standard library context from an async one.
func StdContext(ctx Context) context_.Context {
	return context.Std(ctx)
}
