// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package worker

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/basecomplextech/baselibrary/async"
	"github.com/basecomplextech/baselibrary/status"
	"github.com/stretchr/testify/assert"
)

func testStop[T any](w Worker[T]) {
	select {
	case <-w.Stop():
	case <-time.After(time.Second):
		panic("stop timeout")
	}
}

// Run

func TestWorker_Run__should_process_tasks_in_background(t *testing.T) {
	wg := sync.WaitGroup{}
	total := atomic.Int64{}

	fn := func(ctx async.Context, task int64) status.Status {
		defer wg.Done()

		total.Add(task)
		return status.OK
	}

	w := newWorker(fn)
	w.Start()

	wg.Add(3)
	w.Push(1)
	w.Push(2)
	w.Push(3)

	wg.Wait()
	assert.Equal(t, int64(6), total.Load())
}

func TestWorker_Run__should_stop_on_error(t *testing.T) {
	wg := sync.WaitGroup{}
	total := atomic.Int64{}

	fn := func(ctx async.Context, task int64) status.Status {
		defer wg.Done()

		if total.Add(1) == 3 {
			return status.Error("test error")
		}
		return status.OK
	}

	w := newWorker(fn)
	w.Start()

	wg.Add(3)
	w.Push(1)
	w.Push(2)
	w.Push(3)

	wg.Wait()

	select {
	case <-w.Stopped().Wait():
	case <-time.After(time.Second):
		panic("stop timeout")
	}

	st := w.Status()
	assert.Equal(t, status.Error("test error"), st)
}

// Start

func TestWorker_Start__should_restart_worker(t *testing.T) {
	wg := sync.WaitGroup{}
	total := atomic.Int64{}

	fn := func(ctx async.Context, task int64) status.Status {
		select {
		case <-ctx.Wait():
			return ctx.Status()
		case <-time.After(time.Millisecond * 100):
		}

		total.Add(task)
		wg.Done()
		return status.OK
	}

	w := newWorker(fn)
	w.Start()

	wg.Add(3)
	w.Push(1)
	w.Push(2)
	w.Push(3)

	time.Sleep(50 * time.Millisecond)
	testStop(w)

	w.Start()
	wg.Wait()
	assert.Equal(t, int64(6), total.Load())
}

// Stop

func TestWorker_Stop__should_stop_worker(t *testing.T) {
	fn := func(ctx async.Context, task int64) status.Status {
		return status.OK
	}

	w := newWorker(fn)
	w.Start()

	w.Push(1)
	w.Push(2)
	w.Push(3)

	testStop(w)

	st := w.Status()
	assert.Equal(t, status.Cancelled, st)
}

// Race

func TestWorker__should_not_have_race_conditions(t *testing.T) {
	wg := sync.WaitGroup{}
	done := atomic.Bool{}
	total := atomic.Int64{}

	fn := func(ctx async.Context, task int64) status.Status {
		total.Add(task)
		wg.Done()
		return status.OK
	}

	w := newWorker(fn)
	w.Start()

	taskNum := 20_000
	wg.Add(taskNum)

	for i := 0; i < taskNum; i++ {
		go func() {
			for !w.handling.Load() {
				time.Sleep(time.Millisecond)
			}

			w.Push(1)
		}()
	}

	go func() {
		for !done.Load() {
			w.Start()
			<-w.Stop()
		}
	}()

	wg.Wait()
	done.Store(true)

	assert.Equal(t, int64(taskNum), total.Load())
}
