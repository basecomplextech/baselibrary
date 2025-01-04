// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import (
	"github.com/basecomplextech/baselibrary/async/internal/promise"
	"github.com/basecomplextech/baselibrary/status"
)

// Promise is a future which can be completed.
type Promise[T any] = promise.Promise[T]

// NewPromise returns a pending promise.
func NewPromise[T any]() Promise[T] {
	return promise.New[T]()
}

// Resolved

// Resolved returns a successful future.
func Resolved[T any](result T) Future[T] {
	return promise.Resolved(result)
}

// Rejected returns a rejected future.
func Rejected[T any](st status.Status) Future[T] {
	return promise.Rejected[T](st)
}

// Completed returns a completed future.
func Completed[T any](result T, st status.Status) Future[T] {
	return promise.Completed(result, st)
}
