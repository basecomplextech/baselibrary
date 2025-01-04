// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package pool

import (
	"runtime"
	"sync"
	"testing"
)

const (
	goroutineNum = 100
)

func BenchmarkGC(b *testing.B) {
	p := New()
	wg := &sync.WaitGroup{}

	for i := 0; i < b.N; i++ {
		wg.Add(1)

		p.Go(func() {
			computeRecursive(0, 100)
			wg.Done()
		})
	}

	wg.Wait()

	gcBefore := runtime.NumGoroutine()
	b.ReportMetric(float64(gcBefore), "gc-before")

	runtime.GC()
	runtime.GC()
	runtime.GC()

	gcAfter := runtime.NumGoroutine()
	b.ReportMetric(float64(gcAfter), "gc-after")
}

// Stack 100

func Benchmark_stack100_pool(b *testing.B) {
	p := New()
	stackDepth := 100

	for i := 0; i < b.N; i++ {
		wg := &sync.WaitGroup{}

		for j := 0; j < goroutineNum; j++ {
			wg.Add(1)
			p.Go(func() {
				computeRecursive(0, stackDepth)
				wg.Done()
			})
		}

		wg.Wait()
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N*goroutineNum) / sec
	b.ReportMetric(ops, "ops")
}

func Benchmark_stack100_nopool(b *testing.B) {
	stackDepth := 100

	for i := 0; i < b.N; i++ {
		wg := &sync.WaitGroup{}

		for j := 0; j < goroutineNum; j++ {
			wg.Add(1)
			go func() {
				computeRecursive(0, stackDepth)
				wg.Done()
			}()
		}

		wg.Wait()
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N*goroutineNum) / sec
	b.ReportMetric(ops, "ops")
}

// Stack 1000

func Benchmark_stack1000_pool(b *testing.B) {
	p := New()
	stackDepth := 1000

	for i := 0; i < b.N; i++ {
		wg := &sync.WaitGroup{}

		for j := 0; j < goroutineNum; j++ {
			wg.Add(1)
			p.Go(func() {
				computeRecursive(0, stackDepth)
				wg.Done()
			})
		}

		wg.Wait()
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N*goroutineNum) / sec
	b.ReportMetric(ops, "ops")
}

func Benchmark_stack1000_nopool(b *testing.B) {
	stackDepth := 1000

	for i := 0; i < b.N; i++ {
		wg := &sync.WaitGroup{}

		for j := 0; j < goroutineNum; j++ {
			wg.Add(1)
			go func() {
				computeRecursive(0, stackDepth)
				wg.Done()
			}()
		}

		wg.Wait()
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N*goroutineNum) / sec
	b.ReportMetric(ops, "ops")
}

// single

func Benchmark_stack100_single(b *testing.B) {
	stackDepth := 100

	for i := 0; i < b.N; i++ {
		computeRecursive(0, stackDepth)
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops, "ops")
}

func Benchmark_stack1000_single(b *testing.B) {
	stackDepth := 1000

	for i := 0; i < b.N; i++ {
		computeRecursive(0, stackDepth)
	}

	sec := b.Elapsed().Seconds()
	ops := float64(b.N) / sec
	b.ReportMetric(ops, "ops")
}

// private

func computeRecursive(i int, n int) int {
	for j := 0; j < 100; j++ {
		i++
	}
	if n == 0 {
		return i
	}

	return computeRecursive(i, n-1)
}
