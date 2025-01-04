// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package async

import (
	"github.com/basecomplextech/baselibrary/async/internal/pool"
)

// Pool allows to reuse goroutines with preallocated big stacks.
type Pool = pool.Pool

// NewPool returns a new goroutine pool.
func NewPool() Pool {
	return pool.New()
}
