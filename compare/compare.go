// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package compare

import (
	"github.com/basecomplextech/baselibrary/status"
)

// Compare compares to values and returns the result.
// The result should be 0 if a == b, negative if a < b, and positive if a > b.
type Compare[T any] func(a, b T) int

// CompareError compares to values and returns the result, or an error.
// The result should be 0 if a == b, negative if a < b, and positive if a > b.
type CompareError[T any] func(a, b T) (int, error)

// CompareStatus compares to values and returns the result and a status.
// The result should be 0 if a == b, negative if a < b, and positive if a > b.
type CompareStatus[T any] func(a, b T) (int, status.Status)

// CompareTo

// CompareTo compares a value to another and returns the result.
// The result should be 0 if value == another, <0 if value < another, and >0 if value > another.
//
// Example:
//
//	var another int
//	func(value int) int {
//		return value - another
//	}
type CompareTo[T any] func(value T) int

// CompareToError compares a value to another and returns the result, or an error.
// The result should be 0 if value == another, <0 if value < another, and >0 if value > another.
//
// Example:
//
//	var another int
//	func(value int) (int, error) {
//		return value - another, nil
//	}
type CompareToError[T any] func(value T) (int, error)

// CompareToStatus compares a value to another and returns the result and a status.
// The result should be 0 if value == another, <0 if value < another, and >0 if value > another.
//
// Example:
//
//	var another int
//	func(value int) (int, status.Status) {
//		return value - another, status.OK
//	}
type CompareToStatus[T any] func(value T) (int, status.Status)
