package async

import "github.com/epochtimeout/library/status"

type Result[T any] struct {
	Value  T
	Status status.Status
}

// NewResult returns a new result.
func NewResult[T any](value T, st status.Status) Result[T] {
	return Result[T]{
		Value:  value,
		Status: st,
	}
}
