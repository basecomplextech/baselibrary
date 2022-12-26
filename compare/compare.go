package compare

import (
	"bytes"
	"strings"
	"time"

	"github.com/complex1tech/baselibrary/constraints"
)

type (
	IntFunc     = Func[int]
	Int32Func   = Func[int32]
	Int64Func   = Func[int64]
	UintFunc    = Func[uint]
	Uint32Func  = Func[uint32]
	Uint64Func  = Func[uint64]
	Float32Func = Func[float32]
	Float64Func = Func[float64]
	BytesFunc   = Func[[]byte]
	BinaryFunc  = Func[[]byte]
	StringFunc  = Func[string]
	TimeFunc    = Func[time.Time]
)

// Func is a generic comparison function.
// The result should be 0 if a == b, negative if a < b, and positive if a > b.
type Func[T any] func(a, b T) int

// Less is a generic comparison function.
// The result should be true if a < b, and false otherwise.
type Less[T any] func(a, b T) bool

// Reverse reverses a comparison function.
func Reverse[T any](cmp Func[T]) Func[T] {
	return func(a, b T) int {
		return -cmp(a, b)
	}
}

func Int(a, b int) int {
	return a - b
}

func Int32(a, b int32) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

func Int64(a, b int64) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

func Uint(a, b uint) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

func Uint32(a, b uint32) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

func Uint64(a, b uint64) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

func Float32(a, b float32) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

func Float64(a, b float64) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

func Binary(a, b []byte) int {
	return bytes.Compare(a, b)
}

func Bytes(a, b []byte) int {
	return bytes.Compare(a, b)
}

func String(a, b string) int {
	return strings.Compare(a, b)
}

func Time(a, b time.Time) int {
	switch {
	case a.Equal(b):
		return 0
	case a.Before(b):
		return -1
	default:
		return 1
	}
}

func Duration(a, b time.Duration) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

// Ordered returns a comparison function for a natually ordered type.
func Ordered[T constraints.Ordered]() Func[T] {
	return func(a, b T) int {
		switch {
		case a < b:
			return -1
		case a > b:
			return 1
		default:
			return 0
		}
	}
}
