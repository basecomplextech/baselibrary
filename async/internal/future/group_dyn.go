// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package future

import (
	"github.com/basecomplextech/baselibrary/async/internal/context"
	"github.com/basecomplextech/baselibrary/status"
)

// FutureGroupDyn is a group of futures of different types.
// Use [FutureGroup] for a group of futures of the same type.
type FutureGroupDyn []FutureDyn

// Await waits for the completion of all futures in the group.
// The method returns nil and the context status if the context is cancelled.
func (g FutureGroupDyn) Await(ctx context.Context) status.Status {
	return AwaitAll(ctx, g...)
}

// AwaitAny waits for the completion of any future in the group, and returns its status.
// The method returns -1 and the context status if the context is cancelled.
func (g FutureGroupDyn) AwaitAny(ctx context.Context) (int, status.Status) {
	return AwaitAnyDyn(ctx, g...)
}

// AwaitError waits for any failure of the futures in the group, and returns the error.
// The method returns ok if all futures are successful.
func (g FutureGroupDyn) AwaitError(ctx context.Context) (int, status.Status) {
	return AwaitError(ctx, g...)
}

// Statuses returns the statuses of all futures in the group.
func (g FutureGroupDyn) Statuses() []status.Status {
	result := make([]status.Status, 0, len(g))
	for _, f := range g {
		st := f.Status()
		result = append(result, st)
	}
	return result
}
