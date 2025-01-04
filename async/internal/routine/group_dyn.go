// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package routine

import (
	"github.com/basecomplextech/baselibrary/async/internal"
	"github.com/basecomplextech/baselibrary/async/internal/context"
	"github.com/basecomplextech/baselibrary/async/internal/future"
	"github.com/basecomplextech/baselibrary/status"
)

// RoutineGroupDyn is a group of routines of different types.
// Use [RoutineGroup] for a group of routines of the same type.
type RoutineGroupDyn []RoutineDyn

// Await waits for the completion of all routines in the group.
// The method returns the context status if the context is cancelled.
func (g RoutineGroupDyn) Await(ctx context.Context) status.Status {
	return future.AwaitAll(ctx, g...)
}

// AwaitAny waits for the completion of any routine in the group, and returns its result.
// The method returns -1 and the context status if the context is cancelled.
func (g RoutineGroupDyn) AwaitAny(ctx context.Context) (int, status.Status) {
	return future.AwaitAnyDyn(ctx, g...)
}

// AwaitError waits for any failure of the routines in the group, and returns the error.
// The method returns ok if all routines are successful.
func (g RoutineGroupDyn) AwaitError(ctx context.Context) (int, status.Status) {
	return future.AwaitError(ctx, g...)
}

// Statuses returns the statuses of all routines in the group.
func (g RoutineGroupDyn) Statuses() []status.Status {
	result := make([]status.Status, 0, len(g))
	for _, r := range g {
		st := r.Status()
		result = append(result, st)
	}
	return result
}

// Stop stops all routines in the group.
func (g RoutineGroupDyn) Stop() {
	internal.StopAll(g...)
}

// StopWait stops and awaits all routines in the group.
func (g RoutineGroupDyn) StopWait() {
	internal.StopWaitAll(g...)
}
