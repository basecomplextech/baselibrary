package logging

import (
	"os"
	"testing"
	"time"
)

func BenchmarkLogger(b *testing.B) {
	f, err := os.OpenFile("test.log", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(f.Name())
	defer f.Close()

	w := newConsoleWriter(LevelDebug, true, f)
	l := newLogger("main", true, w)
	start := time.Now()

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		l.Info("Hello, world", "key", 1, "key1", 2, "key2", 3)
	}

	sec := time.Since(start)
	rps := float64(b.N) / sec.Seconds()
	b.ReportMetric(rps, "rps")
}
