package async

import (
	"runtime"
	"sync"

	"github.com/basecomplextech/baselibrary/internal/hashing"
	"github.com/basecomplextech/baselibrary/pools"
	"github.com/basecomplextech/baselibrary/status"
)

// LockMap holds locks for different keys.
//
// The map is a sharded map, which uses a lock per shard.
// The number of shards is equal to the number of CPU cores.
type LockMap[K comparable] interface {
	// Get returns a key lock, the lock must be freed after use.
	//
	// Usage:
	//
	//	m := NewLockMap[int]()
	//
	//	lock := m.Get(123)
	//	defer lock.Free()
	//
	//	select {
	//	case <-lock.Lock():
	//	case <-time.After(time.Second):
	//		return status.Timeout
	//	case <-ctx.Wait():
	//		return ctx.Status()
	//	}
	//	defer lock.Unlock()
	Get(key K) KeyLock

	// Contains returns true if the key is present.
	//
	// Usually it means that the key is locked, but it is not guaranteed.
	// In the latter case the key is unlocked but is not yet freed.
	Contains(key K) bool

	// Lock returns a locked key, the key must be freed after use.
	//
	// Usage:
	//
	//	m := NewLockMap[int]()
	//
	//	lock, st := m.Lock(ctx, 123)
	//	if !st.OK() {
	//		return st
	//	}
	//	defer lock.Free()
	Lock(ctx Context, key K) (LockedKey, status.Status)

	// LockMap locks the map itself, internally it locks all shards.
	//
	// Usage:
	//
	//	m := NewLockMap[int]()
	//
	//	locks := m.LockMap()
	//	defer locks.Free()
	//
	//	for key := range keys {
	//		ok := locks.Contains(key)
	//		// ...
	//	}
	LockMap() LockedMap[K]
}

// LockedMap is an interface to interact with a map locked in exclusive mode.
type LockedMap[K comparable] interface {
	// Contains returns true if the key is present.
	Contains(key K) bool

	// Range ranges over all keys.
	Range(f func(key K) bool)

	// Free unlocks the map itself, internally it unlocks all shards.
	Free()
}

// KeyLock is a single lock for a key, the lock must be freed after use.
type KeyLock interface {
	// Lock returns a channel receiving from which locks the key.
	Lock() <-chan struct{}

	// Unlock unlocks the key lock.
	Unlock()

	// Free frees the acquired key.
	Free()
}

// LockedKey is a locked key which is unlocked when freed.
type LockedKey interface {
	// Free unlocks and freed the key.
	Free()
}

// NewLockMap returns a new lock map.
func NewLockMap[K comparable]() LockMap[K] {
	return newLockMap[K]()
}

// internal

var _ LockMap[any] = &lockMap[any]{}

type lockMap[K comparable] struct {
	shards []lockShard[K]
}

func newLockMap[K comparable]() *lockMap[K] {
	num := runtime.NumCPU()
	shards := make([]lockShard[K], num)

	for i := range shards {
		shards[i] = newLockShard[K]()
	}

	return &lockMap[K]{
		shards: shards,
	}
}

// Get returns a key key, the lock must be freed after use.
func (m *lockMap[K]) Get(key K) KeyLock {
	// Get lock item
	shard := m.shard(key)
	item := shard.get(key)

	// Return key lock
	return &keyLock[K]{item}
}

// Contains returns true if the key is present.
//
// Usually it means that the key is locked, but it is not guaranteed.
// In the latter case the key is unlocked but is not yet freed.
func (m *lockMap[K]) Contains(key K) bool {
	shard := m.shard(key)
	return shard.contains(key)
}

// Lock returns a locked key, the key must be freed after use.
func (m *lockMap[K]) Lock(ctx Context, key K) (LockedKey, status.Status) {
	// Get lock item
	shard := m.shard(key)
	item := shard.get(key)

	// Free if not locked
	ok := false
	defer func() {
		if !ok {
			item.free()
		}
	}()

	// Try lock
	select {
	case <-item.lock:
	default:
		// Lock or wait
		// Context channel is lazily allocated, so try to postpone calling wait.
		select {
		case <-item.lock:
		case <-ctx.Wait():
			return nil, ctx.Status()
		}
	}

	// Return locked key
	k := &lockedKey[K]{item}
	ok = true
	return k, status.OK
}

// LockMap locks the map itself, internally it locks all shards.
func (m *lockMap[K]) LockMap() LockedMap[K] {
	for i := range m.shards {
		shard := &m.shards[i]
		shard.mu.Lock()
	}

	return &lockedMap[K]{m}
}

// unlockMap unlocks the map itself, internally it unlocks all shards.
func (m *lockMap[K]) unlockMap() {
	for i := range m.shards {
		shard := &m.shards[i]
		shard.mu.Unlock()
	}
}

func (m *lockMap[K]) shard(key K) *lockShard[K] {
	index := hashing.Shard(key, len(m.shards))
	return &m.shards[index]
}

// LockedMap

var _ LockedMap[any] = &lockedMap[any]{}

type lockedMap[K comparable] struct {
	m *lockMap[K]
}

// Contains returns true if the key is present.
func (m *lockedMap[K]) Contains(key K) bool {
	shard := m.m.shard(key)
	return shard.containsLocked(key)
}

// Range ranges over all keys.
func (m *lockedMap[K]) Range(f func(key K) bool) {
	for i := range m.m.shards {
		shard := &m.m.shards[i]
		shard.rangeLocked(f)
	}
}

// Free unlocks the map itself, internally it unlocks all shards.
func (m *lockedMap[K]) Free() {
	m.m.unlockMap()
}

// KeyLock

var _ KeyLock = &keyLock[any]{}

type keyLock[K comparable] struct {
	item *lockItem[K]
}

// Lock returns a channel receiving from which locks the key.
func (l *keyLock[K]) Lock() <-chan struct{} {
	return l.item.lock
}

// Unlock unlocks the key lock.
func (l *keyLock[K]) Unlock() {
	l.item.unlock()
}

// Free frees the acquired key.
func (l *keyLock[K]) Free() {
	if l.item == nil {
		panic("free of freed key lock")
	}

	item := l.item
	l.item = nil
	item.free()
}

// LockedKey

var _ LockedKey = &lockedKey[any]{}

type lockedKey[K comparable] struct {
	item *lockItem[K]
}

func (l *lockedKey[K]) Free() {
	if l.item == nil {
		panic("free of freed locked key")
	}

	item := l.item
	l.item = nil

	item.unlock()
	item.free()
}

// shard

type lockShard[K comparable] struct {
	mu    sync.RWMutex
	items map[K]*lockItem[K]
	pool  pools.Pool[*lockItem[K]]
	_     [208]byte // cache line padding
}

func newLockShard[K comparable]() lockShard[K] {
	return lockShard[K]{
		items: make(map[K]*lockItem[K]),
		pool:  pools.NewPool[*lockItem[K]](),
	}
}

func (s *lockShard[K]) get(key K) *lockItem[K] {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Try to get item
	m, ok := s.items[key]
	if ok {
		m.refs++
		return m
	}

	// Make new item with 1 refs
	m, ok = s.pool.Get()
	if ok {
		m.shard = s
		m.key = key
		m.refs = 1
	} else {
		m = newLockItem(s, key)
	}

	// Store item
	s.items[key] = m
	return m
}

func (s *lockShard[K]) free(m *lockItem[K]) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Decrement refs
	if m.refs <= 0 {
		panic("free of freed key lock")
	}
	m.refs--
	if m.refs > 0 {
		return
	}

	// Delete item with 0 refs
	delete(s.items, m.key)

	// Release it to the pool
	m.reset()
	s.pool.Put(m)
}

func (s *lockShard[K]) contains(key K) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.items[key]
	return ok
}

func (s *lockShard[K]) containsLocked(key K) bool {
	_, ok := s.items[key]
	return ok
}

func (s *lockShard[K]) rangeLocked(f func(key K) bool) {
	for key := range s.items {
		if !f(key) {
			break
		}
	}
}

// item

type lockItem[K comparable] struct {
	shard *lockShard[K]
	lock  chan struct{}
	refs  int32

	key K
}

func newLockItem[K comparable](s *lockShard[K], key K) *lockItem[K] {
	m := &lockItem[K]{
		shard: s,
		lock:  make(chan struct{}, 1),
		refs:  1,

		key: key,
	}
	m.lock <- struct{}{}
	return m
}

func (m *lockItem[K]) unlock() {
	select {
	case m.lock <- struct{}{}:
	default:
		panic("unlock of unlocked key lock")
	}
}

func (m *lockItem[K]) free() {
	m.shard.free(m)
}

func (m *lockItem[K]) reset() {
	var zero K
	m.shard = nil
	m.key = zero
	m.refs = 0

	select {
	case m.lock <- struct{}{}:
	default:
	}
}
