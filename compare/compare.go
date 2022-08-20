package compare

import (
	"bytes"
	"strings"
	"time"
)

// Compare is a generic comparison function.
// The result should be 0 if a == b, negative if a < b, and positive if a > b.
type Compare[T any] func(a, b T) int

// Less is a generic comparison function.
// The result should be true if a < b, and false otherwise.
type Less[T any] func(a, b T) bool

// Reverse reverses a comparison function.
func Reverse[T any](cmp Compare[T]) Compare[T] {
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
