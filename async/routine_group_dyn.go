// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import (
	"github.com/basecomplextech/baselibrary/status"
)

// RoutineGroupDyn is a group of routines of different types.
// Use [RoutineGroup] for a group of routines of the same type.
type RoutineGroupDyn []RoutineDyn

// Await waits for the completion of all routines in the group.
// The method returns the context status if the context is cancelled.
func (g RoutineGroupDyn) Await(ctx Context) status.Status {
	return AwaitAll(ctx, g...)
}

// AwaitAny waits for the completion of any routine in the group, and returns its result.
// The method returns -1 and the context status if the context is cancelled.
func (g RoutineGroupDyn) AwaitAny(ctx Context) (int, status.Status) {
	return AwaitAnyDyn(ctx, g...)
}

// AwaitError waits for any failure of the routines in the group, and returns the error.
// The method returns ok if all routines are successful.
func (g RoutineGroupDyn) AwaitError(ctx Context) (int, status.Status) {
	return AwaitError(ctx, g...)
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
	StopAll(g...)
}

// StopWait stops and awaits all routines in the group.
func (g RoutineGroupDyn) StopWait() {
	StopWaitAll(g...)
}
