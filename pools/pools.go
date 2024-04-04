package pools

import (
	"maps"
	"sync"
	"sync/atomic"
)

// Pools is a map of pools by type, experimental.
type Pools struct {
	wmu sync.Mutex
	cur atomic.Pointer[version]
}

// New returns a new pools instance.
func New() *Pools {
	return &Pools{}
}

// Acquire returns a value from a generic pool.
func Acquire[K any, T any](p *Pools) (zero T) {
	pool := Get[K](p)
	v := pool.Get()
	if v == nil {
		return zero
	}
	return v.(T)
}

// Release returns a value to a generic pool.
func Release[K any, T any](p *Pools, v T) {
	pool := Get[K](p)
	pool.Put(v)
}

// Get returns a pool for a type.
func Get[K any](p *Pools) *sync.Pool {
	var key K

	// Fast path
	v := p.cur.Load()
	if v != nil {
		pool, ok := v.m[key]
		if ok {
			return pool
		}
	}

	// Slow path
	p.wmu.Lock()
	defer p.wmu.Unlock()

	// Check again
	v = p.cur.Load()
	if v != nil {
		pool, ok := v.m[key]
		if ok {
			return pool
		}
	}

	// Make or clone version
	var v1 *version
	if v == nil {
		v1 = newVersion()
	} else {
		v1 = v.clone()
	}

	// Add pool, replace version
	pool := &sync.Pool{}
	v1.m[key] = pool
	p.cur.Store(v1)
	return pool
}

// private

type version struct {
	m map[any]*sync.Pool
}

func newVersion() *version {
	return &version{
		m: make(map[any]*sync.Pool),
	}
}

func (v *version) clone() *version {
	m1 := maps.Clone(v.m)
	return &version{m: m1}
}
