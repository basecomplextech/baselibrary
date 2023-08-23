package bufqueue

import (
	"bytes"
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
	waitWrite := 0
	waitRead := 0

	go func() {
		for i := 0; i < b.N; {
			ok, st := q.Write(msg0)
			if !st.OK() {
				b.Fatal(st)
			}

			if ok {
				i++
				continue
			}

			<-q.WaitCanWrite(len(msg0))
			waitWrite++
		}
	}()

	for i := 0; i < b.N; {
		msg1, ok, st := q.Read()
		if !st.OK() {
			b.Fatal(st)
		}

		if ok {
			if !bytes.Equal(msg0, msg1) {
				b.Fatal("invalid message")
			}
			i++
			continue
		}

		<-q.Wait()
		waitRead++
	}

	sec := time.Since(t0).Seconds()
	ops := float64(b.N) / float64(sec)

	b.ReportMetric(ops/1000_000, "mops")
	b.ReportMetric(float64(q.maxCap), "max-cap")
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
	waitWrite := 0
	waitRead := 0

	go func() {
		for i := 0; i < b.N; {
			ok, st := q.Write(msg0)
			if !st.OK() {
				b.Fatal(st)
			}

			if ok {
				i++
				continue
			}

			<-q.WaitCanWrite(len(msg0))
			waitWrite++
		}
	}()

	for i := 0; i < b.N; {
		msg1, ok, st := q.Read()
		if !st.OK() {
			b.Fatal(st)
		}

		if ok {
			if !bytes.Equal(msg0, msg1) {
				b.Fatal("invalid message")
			}
			i++
			continue
		}

		<-q.Wait()
		waitRead++
	}

	sec := time.Since(t0).Seconds()
	ops := float64(b.N) / float64(sec)

	b.ReportMetric(ops/1000_000, "mops")
	b.ReportMetric(float64(q.maxCap), "max-cap")
	b.ReportMetric(float64(waitWrite), "wait-write")
	b.ReportMetric(float64(waitRead), "wait-read")
}

func BenchmarkQueue_1024b(b *testing.B) {
	h := heap.New()
	q := newQueue(h, 128*1024)
	msg0 := bytes.Repeat([]byte("a"), 1024-4)

	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(int64(len(msg0)))

	t0 := time.Now()
	waitWrite := 0
	waitRead := 0

	go func() {
		for i := 0; i < b.N; {
			ok, st := q.Write(msg0)
			if !st.OK() {
				b.Fatal(st)
			}

			if ok {
				i++
				continue
			}

			<-q.WaitCanWrite(len(msg0))
			waitWrite++
		}
	}()

	for i := 0; i < b.N; {
		msg1, ok, st := q.Read()
		if !st.OK() {
			b.Fatal(st)
		}

		if ok {
			if !bytes.Equal(msg0, msg1) {
				b.Fatal("invalid message")
			}
			i++
			continue
		}

		<-q.Wait()
		waitRead++
	}

	sec := time.Since(t0).Seconds()
	ops := float64(b.N) / float64(sec)

	b.ReportMetric(ops/1000_000, "mops")
	b.ReportMetric(float64(q.maxCap), "max-cap")
	b.ReportMetric(float64(waitWrite), "wait-write")
	b.ReportMetric(float64(waitRead), "wait-read")
}
