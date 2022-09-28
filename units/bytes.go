package units

import (
	"fmt"
	"math"
)

const (
	_ Bytes = 1 << (10 * iota)
	KiB
	MiB
	GiB
	TiB
	PiB
	EiB
)

// Bytes formats a byte count into a human readable string.
// TODO: Test bytes, including negative bytes.
type Bytes int64

func (b Bytes) String() string {
	return ByteString(int64(b))
}

// ByteString produces a human readable representation of a byte count.
//
// ByteString(82854982) -> 79mb
func ByteString(s int64) string {
	sizes := []string{"b", "kb", "mb", "gb", "tb", "pb", "eb"}
	return byteString(s, 1024, sizes)
}

func byteString(s int64, base float64, sizes []string) string {
	if s < 10 {
		return fmt.Sprintf("%db", s)
	}

	e := math.Floor(logn(float64(s), base))
	var suffix string
	if e < float64(len(sizes)) {
		suffix = sizes[int(e)]
	} else {
		suffix = sizes[len(sizes)-1]
	}

	val := math.Floor(float64(s)/math.Pow(base, e)*10+0.5) / 10
	f := "%.0f%s"
	if val < 10 {
		f = "%.1f%s"
	}

	return fmt.Sprintf(f, val, suffix)
}

func logn(n, b float64) float64 {
	return math.Log(n) / math.Log(b)
}
