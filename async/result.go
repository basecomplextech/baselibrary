// Copyright 2022 Ivan Korobkov. All rights reserved.

package async

import "github.com/basecomplextech/baselibrary/status"

// Result is an async result which combines a value and a status.
type Result[T any] struct {
	Value  T
	Status status.Status
}

// Unwrap returns the value and the status.
func (r Result[T]) Unwrap() (T, status.Status) {
	return r.Value, r.Status
}
