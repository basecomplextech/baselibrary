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

type (
	// Func is a function which returns the result.
	Func[T any] = func(ctx Context) (T, status.Status)

	// Func1 is a single argument function which returns the result.
	Func1[T any, A any] = func(ctx Context, arg A) (T, status.Status)

	// FuncVoid is a function which returns no result.
	FuncVoid = func(ctx Context) status.Status

	// FuncVoid1 is a single argument function which returns no result
	FuncVoid1[A any] = func(ctx Context, arg A) status.Status
)

// New

// NewRoutine returns a new routine, but does not start it.
func NewRoutine[T any](fn Func[T]) Routine[T] {
	return routine.New(fn)
}

// NewRoutineVoid returns a new routine, but does not start it.
func NewRoutineVoid(fn FuncVoid) RoutineVoid {
	return routine.NewVoid(fn)
}

// Run

// Run runs a function in a new routine, and returns the result, recovers on panics.
func Run[T any](fn Func[T]) Routine[T] {
	return routine.Run(fn)
}

// Run1 runs a function in a new routine, and returns the result, recovers on panics.
func Run1[T any, A any](fn Func1[T, A], arg A) Routine[T] {
	return routine.Run1(fn, arg)
}

// RunVoid

// RunVoid runs a procedure in a new routine, recovers on panics.
func RunVoid(fn FuncVoid) RoutineVoid {
	return routine.RunVoid(fn)
}

// RunVoid1 runs a procedure in a new routine, recovers on panics.
func RunVoid1[A any](fn FuncVoid1[A], arg A) RoutineVoid {
	return routine.RunVoid1(fn, arg)
}
