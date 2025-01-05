// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import (
	"github.com/basecomplextech/baselibrary/async/internal/routine"
	"github.com/basecomplextech/baselibrary/status"
)

type (
	// Routine is an async routine which returns the result as a future, recovers on panics,
	// and can be cancelled.
	Routine[T any] = routine.Routine[T]

	// RoutineDyn is a routine interface without generics, i.e. Routine[?].
	RoutineDyn = routine.RoutineDyn

	// RoutineVoid is a routine which has no result.
	RoutineVoid = routine.RoutineVoid
)

type (
	// RoutineGroup is a group of routines of the same type.
	// Use [RoutineGroupDyn] for a group of routines of different types.
	RoutineGroup[T any] = routine.RoutineGroup[T]

	// RoutineGroupDyn is a group of routines of different types.
	// Use [RoutineGroup] for a group of routines of the same type.
	RoutineGroupDyn = routine.RoutineGroupDyn
)

// Go

// Go runs a function in a new routine, recovers on panics.
func Go(fn func(ctx Context) status.Status) RoutineVoid {
	return routine.Go(fn)
}

// Run runs a function in a new routine, and returns the result, recovers on panics.
func Run[T any](fn func(ctx Context) (T, status.Status)) Routine[T] {
	return routine.Run(fn)
}

// Exited returns a routine which has exited with the given result and status.
func Exited[T any](result T, st status.Status) Routine[T] {
	return routine.Exited(result, st)
}
