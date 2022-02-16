package async

type Result[T any] struct {
	Value T
	Err   error
}

// NewResult returns a new result.
func NewResult[T any](value T, err error) Result[T] {
	return Result[T]{
		Value: value,
		Err:   err,
	}
}
