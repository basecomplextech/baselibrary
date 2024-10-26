// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package routinepool

import (
	"testing"
	"time"
)

func TestPool_Go__should_run_function_in_the_pool(t *testing.T) {
	p := New()

	done := make(chan struct{})
	p.Go(func() {
		close(done)
	})

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("timeout")
	}
}

func TestPool_Go__should_reuse_workers(t *testing.T) {
	p := New()

	for i := 0; i < 10; i++ {
		done := make(chan struct{})
		p.Go(func() {
			close(done)
		})

		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("timeout")
		}
	}
}
