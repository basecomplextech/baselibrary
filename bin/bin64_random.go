package bin

import (
	"crypto/rand"
	"encoding/binary"
	"time"
)

// Random64 returns a random bin64.
func Random64() Bin64 {
	u := Bin64{}
	if _, err := rand.Read(u[:]); err != nil {
		panic(err)
	}
	return u
}

// TimeRandom64 returns a time-random bin64 with a second resolution.
func TimeRandom64() Bin64 {
	u := Bin64{}

	now := time.Now()
	ts := now.UnixNano() / int64(time.Second)
	binary.BigEndian.PutUint32(u[:], uint32(ts))

	if _, err := rand.Read(u[4:]); err != nil {
		panic(err)
	}
	return u
}
