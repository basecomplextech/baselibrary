// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package context

import (
	"sync"
	"time"

	"github.com/basecomplextech/baselibrary/opt"
)

// timer is guarded with a mutex to prevent data race in constructor with immediate timeout.
type timer struct {
	mu    sync.Mutex
	timer opt.Opt[*time.Timer]
}

func (t *timer) set(timer *time.Timer) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.timer = opt.New(timer)
}

func (t *timer) stop() {
	t.mu.Lock()
	defer t.mu.Unlock()

	timer, ok := t.timer.Clear()
	if !ok {
		return
	}

	timer.Stop()
}
