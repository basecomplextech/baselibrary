// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package compare

import (
	"github.com/basecomplextech/baselibrary/status"
)

type (
	// Compare compares to values and returns the result.
	// The result should be 0 if a == b, negative if a < b, and positive if a > b.
	Compare[T any] func(a, b T) int

	// CompareError compares to values and returns the result, or an error.
	// The result should be 0 if a == b, negative if a < b, and positive if a > b.
	CompareError[T any] func(a, b T) (int, error)

	// CompareStatus compares to values and returns the result and a status.
	// The result should be 0 if a == b, negative if a < b, and positive if a > b.
	CompareStatus[T any] func(a, b T) (int, status.Status)
)

// CompareTo

type (
	// CompareTo compares a value to another and returns the result.
	// The result should be 0 if value == another, <0 if value < another, and >0 if value > another.
	//
	// Example:
	//
	//	var another int
	//	func(value int) int {
	//		return value - another
	//	}
	CompareTo[T any] func(value T) int

	// CompareToError compares a value to another and returns the result, or an error.
	// The result should be 0 if value == another, <0 if value < another, and >0 if value > another.
	//
	// Example:
	//
	//	var another int
	//	func(value int) (int, error) {
	//		return value - another, nil
	//	}
	CompareToError[T any] func(value T) (int, error)

	// CompareToStatus compares a value to another and returns the result and a status.
	// The result should be 0 if value == another, <0 if value < another, and >0 if value > another.
	//
	// Example:
	//
	//	var another int
	//	func(value int) (int, status.Status) {
	//		return value - another, status.OK
	//	}
	CompareToStatus[T any] func(value T) (int, status.Status)
)

// Aliases

type (
	CompareBytes       = Compare[[]byte]
	CompareBytesError  = CompareError[[]byte]
	CompareBytesStatus = CompareStatus[[]byte]
)
