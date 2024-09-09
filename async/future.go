// Copyright 2021 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import (
	"github.com/basecomplextech/baselibrary/status"
)

// Future represents a result available in the future.
type Future[T any] interface {
	// Wait returns a channel which is closed when the future is complete.
	Wait() <-chan struct{}

	// Result returns a value and a status.
	Result() (T, status.Status)

	// Status returns a status or none.
	Status() status.Status
}

// FutureDyn is a future interface without generics, i.e. Future[?].
type FutureDyn interface {
	// Wait returns a channel which is closed when the future is complete.
	Wait() <-chan struct{}

	// Status returns a status or none.
	Status() status.Status
}

// Constructors

// Resolved returns a successful future.
func Resolved[T any](result T) Future[T] {
	p := newPromise[T]()
	p.Complete(result, status.OK)
	return p
}

// Rejected returns a rejected future.
func Rejected[T any](st status.Status) Future[T] {
	var zero T
	p := newPromise[T]()
	p.Complete(zero, st)
	return p
}

// Completed returns a completed future.
func Completed[T any](result T, st status.Status) Future[T] {
	p := newPromise[T]()
	p.Complete(result, st)
	return p
}
