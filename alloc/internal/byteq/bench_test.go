// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package byteq

import (
	"bytes"
	"sync/atomic"
	"testing"
	"time"

	"github.com/basecomplextech/baselibrary/alloc/internal/heap"
)

func BenchmarkQueue_16b(b *testing.B) {
	h := heap.New()
	q := newQueue(h, 128*1024)
	msg0 := []byte("hello, world") // 12+4

	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(int64(len(msg0)))

	t0 := time.Now()
	waitRead := 0
	waitWrite := 0

	N := b.N
	go func() {
		defer q.Close()

		for i := 0; i < N; {
			ok, st := q.Write(msg0)
			if !st.OK() {
				b.Fatal(st)
			}

			if ok {
				i++
				continue
			}

			<-q.WriteWait(len(msg0))
			waitWrite++
		}
	}()

	for {
		msg1, ok, st := q.Read()
		if !st.OK() {
			break
		}

		if ok {
			if !bytes.Equal(msg0, msg1) {
				b.Fatal("invalid message")
			}
			continue
		}

		<-q.ReadWait()
		waitRead++
	}

	sec := time.Since(t0).Seconds()
	ops := float64(b.N) / float64(sec)

	b.ReportMetric(ops/1000_000, "mops")
	b.ReportMetric(float64(waitWrite), "wait-write")
	b.ReportMetric(float64(waitRead), "wait-read")
}

func BenchmarkQueue_128b(b *testing.B) {
	h := heap.New()
	q := newQueue(h, 128*1024)
	msg0 := bytes.Repeat([]byte("a"), 128-4)

	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(int64(len(msg0)))

	t0 := time.Now()
	waitRead := 0
	waitWrite := 0

	go func() {
		defer q.Close()

		for i := 0; i < b.N; {
			ok, st := q.Write(msg0)
			if !st.OK() {
				b.Fatal(st)
			}

			if ok {
				i++
				continue
			}

			<-q.WriteWait(len(msg0))
			waitWrite++
		}
	}()

	for {
		msg1, ok, st := q.Read()
		if !st.OK() {
			break
		}

		if ok {
			if !bytes.Equal(msg0, msg1) {
				b.Fatal("invalid message")
			}
			continue
		}

		<-q.ReadWait()
		waitRead++
	}

	sec := time.Since(t0).Seconds()
	ops := float64(b.N) / float64(sec)

	b.ReportMetric(ops/1000_000, "mops")
	b.ReportMetric(float64(waitWrite), "wait-write")
	b.ReportMetric(float64(waitRead), "wait-read")
}

func BenchmarkQueue_1kb(b *testing.B) {
	h := heap.New()
	q := newQueue(h, 1024*1024)
	msg0 := bytes.Repeat([]byte("a"), 1024-4)

	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(int64(len(msg0)))

	t0 := time.Now()
	waitRead := 0
	waitWrite := 0

	go func() {
		defer q.Close()

		for i := 0; i < b.N; {
			ok, st := q.Write(msg0)
			if !st.OK() {
				b.Fatal(st)
			}

			if ok {
				i++
				continue
			}

			<-q.WriteWait(len(msg0))
			waitWrite++
		}
	}()

	for {
		msg1, ok, st := q.Read()
		if !st.OK() {
			break
		}

		if ok {
			if !bytes.Equal(msg0, msg1) {
				b.Fatal("invalid message")
			}
			continue
		}

		<-q.ReadWait()
		waitRead++
	}

	sec := time.Since(t0).Seconds()
	ops := float64(b.N) / float64(sec)

	b.ReportMetric(ops/1000_000, "mops")
	b.ReportMetric(float64(waitWrite), "wait-write")
	b.ReportMetric(float64(waitRead), "wait-read")
}

// Parallel

func BenchmarkQueue_16b_Parallel(b *testing.B) {
	h := heap.New()
	q := newQueue(h, 128*1024)
	msg0 := []byte("hello, world") // 12+4

	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(int64(len(msg0)))

	t0 := time.Now()
	waitRead := 0
	waitWrite := int64(0)

	done := make(chan struct{})
	go func() {
		defer close(done)

		for {
			msg1, ok, st := q.Read()
			if !st.OK() {
				break
			}

			if ok {
				if !bytes.Equal(msg0, msg1) {
					b.Fatal("invalid message")
				}
				continue
			}

			<-q.ReadWait()
			waitRead++
		}
	}()

	b.RunParallel(func(p *testing.PB) {
	outer:
		for p.Next() {
			// Retry in loop, not in p.Next(), otherwise, we'll block forever.
			for {
				ok, st := q.Write(msg0)
				if !st.OK() {
					b.Fatal(st)
				}
				if ok {
					continue outer
				}

				<-q.WriteWait(len(msg0))
				atomic.AddInt64(&waitWrite, 1)
			}
		}
	})

	q.Close()
	<-done

	sec := time.Since(t0).Seconds()
	ops := float64(b.N) / float64(sec)

	b.ReportMetric(ops/1000_000, "mops")
	b.ReportMetric(float64(waitWrite), "wait-write")
	b.ReportMetric(float64(waitRead), "wait-read")
}

func BenchmarkQueue_128kb_Parallel(b *testing.B) {
	h := heap.New()
	q := newQueue(h, 1024*1024)
	msg0 := bytes.Repeat([]byte("a"), 128*1024)

	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(int64(len(msg0)))

	t0 := time.Now()
	waitRead := 0
	waitWrite := int64(0)

	done := make(chan struct{})
	go func() {
		defer close(done)

		for {
			msg1, ok, st := q.Read()
			if !st.OK() {
				break
			}

			if ok {
				if !bytes.Equal(msg0, msg1) {
					b.Fatal("invalid message")
				}
				continue
			}

			<-q.ReadWait()
			waitRead++
		}
	}()

	b.RunParallel(func(p *testing.PB) {
	outer:
		for p.Next() {
			// Retry in loop, not in p.Next(), otherwise, we'll block forever.
			for {
				ok, st := q.Write(msg0)
				if !st.OK() {
					b.Fatal(st)
				}
				if ok {
					continue outer
				}

				<-q.WriteWait(len(msg0))
				atomic.AddInt64(&waitWrite, 1)
			}
		}
	})

	q.Close()
	<-done

	sec := time.Since(t0).Seconds()
	ops := float64(b.N) / float64(sec)

	b.ReportMetric(ops/1000_000, "mops")
	b.ReportMetric(float64(waitWrite), "wait-write")
	b.ReportMetric(float64(waitRead), "wait-read")
}
