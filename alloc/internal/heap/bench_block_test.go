// Copyright 2025 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package heap

import (
	"testing"
	"time"
)

func BenchmarkBlock_Alloc(b *testing.B) {
	h := New()

	b.ResetTimer()
	b.ReportAllocs()

	t0 := time.Now()
	size := 8
	block := h.Alloc(size)

	for i := 0; i < b.N; i++ {
		ptr := block.Alloc(size)
		if ptr == nil {
			block.Reset()
		}
	}

	sec := time.Since(t0).Seconds()
	ops := float64(b.N) / float64(sec)

	b.SetBytes(int64(size))
	b.ReportMetric(ops/1000_000, "mops")
}
