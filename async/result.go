package async

type Result[T any] struct {
	Err   error
	Value T
}

// NewResult returns a new result.
func NewResult[T any](value T, err error) Result[T] {
	return Result[T]{
		Err:   err,
		Value: value,
	}
}
