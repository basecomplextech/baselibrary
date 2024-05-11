package pools

import (
	"maps"
	"reflect"
	"sync"
	"sync/atomic"
)

// Pools is a map of pools by type, experimental.
type Pools struct {
	wmu sync.Mutex
	cur atomic.Pointer[version]
}

// NewPools returns a new pools instance.
func NewPools() *Pools {
	return &Pools{}
}

// Values

// Acquire returns a value from a generic pool.
func Acquire[T any](p *Pools) (T, bool) {
	pool := GetPool[T](p)
	return pool.Get()
}

// Acquire1 returns a value from a generic pool, and its pool.
func Acquire1[T any](p *Pools) (T, bool, Pool[T]) {
	pool := GetPool[T](p)
	v, ok := pool.Get()
	return v, ok, pool
}

// Release returns a value to a generic pool.
func Release[T any](p *Pools, v T) {
	pool := GetPool[T](p)
	pool.Put(v)
}

// Pool

// GetPool returns a pool for a type.
func GetPool[T any](p *Pools) Pool[T] {
	key := reflect.TypeFor[T]()

	// Fast path
	v := p.cur.Load()
	if v != nil {
		pool, ok := v.m[key]
		if ok {
			return pool.(Pool[T])
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
			return pool.(Pool[T])
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
	pool := newPool[T](nil)
	v1.m[key] = pool
	p.cur.Store(v1)
	return pool
}

// private

type version struct {
	m map[reflect.Type]any // Pool[T]
}

func newVersion() *version {
	return &version{
		m: make(map[reflect.Type]any),
	}
}

func (v *version) clone() *version {
	m1 := maps.Clone(v.m)
	return &version{m: m1}
}
