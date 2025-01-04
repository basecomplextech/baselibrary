// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import (
	"github.com/basecomplextech/baselibrary/async/internal/routinepool"
)

// RoutinePool allows to reuse goroutines with preallocated big stacks.
type RoutinePool = routinepool.RoutinePool

// NewRoutinePool returns a new goroutine pool.
//
// Use [GoPool] and [RunPool] to run functions in the pool.
func NewRoutinePool() RoutinePool {
	return routinepool.New()
}
