// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import "github.com/basecomplextech/baselibrary/async/internal/future"

type (
	// Future represents a result available in the future.
	Future[T any] = future.Future[T]

	// FutureDyn is a future interface without generics, i.e. Future[?].
	FutureDyn = future.FutureDyn
)

type (
	// FutureGroup is a group of futures of the same type.
	// Use [FutureGroupDyn] for a group of futures of different types.
	FutureGroup[T any] = future.FutureGroup[T]

	// FutureGroupDyn is a group of futures of different types.
	// Use [FutureGroup] for a group of futures of the same type.
	FutureGroupDyn = future.FutureGroupDyn
)

type (
	// Result is an async result which combines a value and a status.
	Result[T any] = future.Result[T]
)
