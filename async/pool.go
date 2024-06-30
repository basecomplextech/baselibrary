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
