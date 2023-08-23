package msgqueue

import (
	"bytes"
	"testing"
	"time"

	"github.com/basecomplextech/baselibrary/alloc/internal/heap"
)

func BenchmarkQueue(b *testing.B) {
	h := heap.New()
	q := newQueue(h, 128*1024)
	msg0 := []byte("hello, world")

	b.ResetTimer()
	b.ReportAllocs()
	b.SetBytes(int64(len(msg0)))

	t0 := time.Now()

	go func() {
		for i := 0; i < b.N; i++ {
			ok, st := q.Write(msg0)
			switch {
			case !st.OK():
				b.Fatal(st)
			case !ok:
				<-q.WaitNotFull(len(msg0))
			}
		}
	}()

	for i := 0; i < b.N; i++ {
		msg1, ok, st := q.Read()
		switch {
		case !st.OK():
			b.Fatal(st)

		case !ok:
			<-q.Wait()

		case !bytes.Equal(msg0, msg1):
			b.Fatal("invalid message")
		}
	}

	sec := time.Since(t0).Seconds()
	ops := float64(b.N) / float64(sec)
	b.ReportMetric(ops/1000_000, "mops")
	b.ReportMetric(float64(q.maxCap), "maxcap")
}
