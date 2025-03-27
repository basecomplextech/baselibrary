// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package compare

import (
	"bytes"
	"strings"
	"time"

	"github.com/basecomplextech/baselibrary/bin"
	"github.com/basecomplextech/baselibrary/status"
)

// Bool

func Bool(a, b bool) int {
	switch {
	case a == b:
		return 0
	case a:
		return 1
	default:
		return -1
	}
}

func BoolError(a, b bool) (int, error) {
	switch {
	case a == b:
		return 0, nil
	case a:
		return 1, nil
	default:
		return -1, nil
	}
}

func BoolStatus(a, b bool) (int, status.Status) {
	switch {
	case a == b:
		return 0, status.OK
	case a:
		return 1, status.OK
	default:
		return -1, status.OK
	}
}

// Byte

func Byte(a, b byte) int {
	return int(a - b)
}

func ByteError(a, b byte) (int, error) {
	return int(a - b), nil
}

func ByteStatus(a, b byte) (int, status.Status) {
	return int(a - b), status.OK
}

// Int

func Int(a, b int) int {
	return a - b
}

func IntError(a, b int) (int, error) {
	return a - b, nil
}

func IntStatus(a, b int) (int, status.Status) {
	return a - b, status.OK
}

// Int16

func Int16(a, b int16) int {
	return int(a - b)
}

func Int16Error(a, b int16) (int, error) {
	return int(a - b), nil
}

func Int16Status(a, b int16) (int, status.Status) {
	return int(a - b), status.OK
}

// Int32

func Int32(a, b int32) int {
	return int(a - b)
}

func Int32Error(a, b int32) (int, error) {
	return int(a - b), nil
}

func Int32Status(a, b int32) (int, status.Status) {
	return int(a - b), status.OK
}

// Int64

func Int64(a, b int64) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

func Int64Error(a, b int64) (int, error) {
	switch {
	case a < b:
		return -1, nil
	case a > b:
		return 1, nil
	}
	return 0, nil
}

func Int64Status(a, b int64) (int, status.Status) {
	switch {
	case a < b:
		return -1, status.OK
	case a > b:
		return 1, status.OK
	}
	return 0, status.OK
}

// Uint

func Uint(a, b uint) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

func UintError(a, b uint) (int, error) {
	switch {
	case a < b:
		return -1, nil
	case a > b:
		return 1, nil
	}
	return 0, nil
}

func UintStatus(a, b uint) (int, status.Status) {
	switch {
	case a < b:
		return -1, status.OK
	case a > b:
		return 1, status.OK
	}
	return 0, status.OK
}

// Uint16

func Uint16(a, b uint16) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

func Uint16Error(a, b uint16) (int, error) {
	switch {
	case a < b:
		return -1, nil
	case a > b:
		return 1, nil
	}
	return 0, nil
}

func Uint16Status(a, b uint16) (int, status.Status) {
	switch {
	case a < b:
		return -1, status.OK
	case a > b:
		return 1, status.OK
	}
	return 0, status.OK
}

// Uint32

func Uint32(a, b uint32) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

func Uint32Error(a, b uint32) (int, error) {
	switch {
	case a < b:
		return -1, nil
	case a > b:
		return 1, nil
	}
	return 0, nil
}

func Uint32Status(a, b uint32) (int, status.Status) {
	switch {
	case a < b:
		return -1, status.OK
	case a > b:
		return 1, status.OK
	}
	return 0, status.OK
}

// Uint64

func Uint64(a, b uint64) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

func Uint64Error(a, b uint64) (int, error) {
	switch {
	case a < b:
		return -1, nil
	case a > b:
		return 1, nil
	}
	return 0, nil
}

func Uint64Status(a, b uint64) (int, status.Status) {
	switch {
	case a < b:
		return -1, status.OK
	case a > b:
		return 1, status.OK
	}
	return 0, status.OK
}

// Float32

func Float32(a, b float32) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

func Float32Error(a, b float32) (int, error) {
	switch {
	case a < b:
		return -1, nil
	case a > b:
		return 1, nil
	}
	return 0, nil
}

func Float32Status(a, b float32) (int, status.Status) {
	switch {
	case a < b:
		return -1, status.OK
	case a > b:
		return 1, status.OK
	}
	return 0, status.OK
}

// Float64

func Float64(a, b float64) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

func Float64Error(a, b float64) (int, error) {
	switch {
	case a < b:
		return -1, nil
	case a > b:
		return 1, nil
	}
	return 0, nil
}

func Float64Status(a, b float64) (int, status.Status) {
	switch {
	case a < b:
		return -1, status.OK
	case a > b:
		return 1, status.OK
	}
	return 0, status.OK
}

// Bin64

func Bin64(a, b bin.Bin64) int {
	return a.Compare(b)
}

func Bin64Error(a, b bin.Bin64) (int, error) {
	return a.Compare(b), nil
}

func Bin64Status(a, b bin.Bin64) (int, status.Status) {
	return a.Compare(b), status.OK
}

// Bin128

func Bin128(a, b bin.Bin128) int {
	return a.Compare(b)
}

func Bin128Error(a, b bin.Bin128) (int, error) {
	return a.Compare(b), nil
}

func Bin128Status(a, b bin.Bin128) (int, status.Status) {
	return a.Compare(b), status.OK
}

// Bin256

func Bin256(a, b bin.Bin256) int {
	return a.Compare(b)
}

func Bin256Error(a, b bin.Bin256) (int, error) {
	return a.Compare(b), nil
}

func Bin256Status(a, b bin.Bin256) (int, status.Status) {
	return a.Compare(b), status.OK
}

// Bytes

func Bytes(a, b []byte) int {
	return bytes.Compare(a, b)
}

func BytesError(a, b []byte) (int, error) {
	return bytes.Compare(a, b), nil
}

func BytesStatus(a, b []byte) (int, status.Status) {
	return bytes.Compare(a, b), status.OK
}

// String

func String(a, b string) int {
	return strings.Compare(a, b)
}

func StringError(a, b string) (int, error) {
	return strings.Compare(a, b), nil
}

func StringStatus(a, b string) (int, status.Status) {
	return strings.Compare(a, b), status.OK
}

// Time

func Time(a, b time.Time) int {
	switch {
	case a.Equal(b):
		return 0
	case a.Before(b):
		return -1
	}
	return 1
}

func TimeError(a, b time.Time) (int, error) {
	switch {
	case a.Equal(b):
		return 0, nil
	case a.Before(b):
		return -1, nil
	}
	return 1, nil
}

func TimeStatus(a, b time.Time) (int, status.Status) {
	switch {
	case a.Equal(b):
		return 0, status.OK
	case a.Before(b):
		return -1, status.OK
	}
	return 1, status.OK
}

// Duration

func Duration(a, b time.Duration) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	}
	return 0
}

func DurationError(a, b time.Duration) (int, error) {
	switch {
	case a < b:
		return -1, nil
	case a > b:
		return 1, nil
	}
	return 0, nil
}

func DurationStatus(a, b time.Duration) (int, status.Status) {
	switch {
	case a < b:
		return -1, status.OK
	case a > b:
		return 1, status.OK
	}
	return 0, status.OK
}
