// Copyright 2022 Ivan Korobkov. All rights reserved.

package compactint

import (
	"encoding/binary"
)

// Size

// ReverseSize decodes and returns the size of an int from the b end, or 0 on error.
func ReverseSize(b []byte) int {
	if len(b) == 0 {
		return 0
	}

	f := b[len(b)-1]
	switch f {
	default:
		return 1

	case 0xfd:
		if len(b) < 3 {
			return 0
		}
		return 3

	case 0xfe:
		if len(b) < 5 {
			return 0
		}
		return 5

	case 0xff:
		if len(b) < 9 {
			return 0
		}
		return 9
	}
}

// Signed

// ReverseInt32 decodes an int32 value from the b end and the number of bytes read.
func ReverseInt32(b []byte) (int32, int) {
	if len(b) == 0 {
		return 0, 0
	}

	var ux uint32
	var size int

	f := b[len(b)-1]
	switch f {
	default:
		ux, size = uint32(f), 1

	case 0xfd:
		if len(b) < 3 {
			return 0, 0
		}

		off := len(b) - 3
		v := binary.BigEndian.Uint16(b[off : off+2])
		ux, size = uint32(v), 3

	case 0xfe:
		if len(b) < 5 {
			return 0, 0
		}

		off := len(b) - 5
		v := binary.BigEndian.Uint32(b[off : off+4])
		ux, size = uint32(v), 5

	case 0xff:
		return 0, -1
	}

	// ok to continue in presence of error
	x := int32(ux >> 1)
	if ux&1 != 0 {
		x = ^x
	}
	return x, size
}

// ReverseInt64 decodes an int64 value from the b end and the number of bytes read.
func ReverseInt64(b []byte) (int64, int) {
	if len(b) == 0 {
		return 0, 0
	}

	var ux uint64
	var size int

	f := b[len(b)-1]
	switch f {
	default:
		ux, size = uint64(f), 1

	case 0xfd:
		if len(b) < 3 {
			return 0, 0
		}

		off := len(b) - 3
		v := binary.BigEndian.Uint16(b[off : off+2])
		ux, size = uint64(v), 3

	case 0xfe:
		if len(b) < 5 {
			return 0, 0
		}

		off := len(b) - 5
		v := binary.BigEndian.Uint32(b[off : off+4])
		ux, size = uint64(v), 5

	case 0xff:
		if len(b) < 9 {
			return 0, 0
		}

		off := len(b) - 9
		v := binary.BigEndian.Uint64(b[off : off+8])
		ux, size = v, 9
	}

	// ok to continue in presence of error
	x := int64(ux >> 1)
	if ux&1 != 0 {
		x = ^x
	}
	return x, size
}

// PutReverseInt32 encodes an int32 into the b end and returns the number of bytes written.
func PutReverseInt32(buf []byte, x int32) int {
	ux := uint32(x) << 1
	if x < 0 {
		ux = ^ux
	}
	return PutReverseUint32(buf, ux)
}

// PutReverseInt64 encodes an int64 into the b end and returns the number of bytes written.
func PutReverseInt64(buf []byte, x int64) int {
	ux := uint64(x) << 1
	if x < 0 {
		ux = ^ux
	}
	return PutReverseUint64(buf, ux)
}

// Unsigned

// ReverseUint32 returns a uint32 value and the number of bytes read.
func ReverseUint32(b []byte) (uint32, int) {
	if len(b) == 0 {
		return 0, 0
	}

	f := b[len(b)-1]
	switch f {
	default:
		return uint32(f), 1

	case 0xfd:
		if len(b) < 3 {
			return 0, 0
		}

		off := len(b) - 3
		v := binary.BigEndian.Uint16(b[off : off+2])
		return uint32(v), 3

	case 0xfe:
		if len(b) < 5 {
			return 0, 0
		}

		off := len(b) - 5
		v := binary.BigEndian.Uint32(b[off : off+4])
		return uint32(v), 5

	case 0xff:
		return 0, -1
	}
}

// ReverseUint64 returns a uint64 value and the number of bytes read.
func ReverseUint64(b []byte) (uint64, int) {
	if len(b) == 0 {
		return 0, 0
	}

	f := b[len(b)-1]
	switch f {
	default:
		return uint64(f), 1

	case 0xfd:
		if len(b) < 3 {
			return 0, 0
		}

		off := len(b) - 3
		v := binary.BigEndian.Uint16(b[off : off+2])
		return uint64(v), 3

	case 0xfe:
		if len(b) < 5 {
			return 0, 0
		}

		off := len(b) - 5
		v := binary.BigEndian.Uint32(b[off : off+4])
		return uint64(v), 5

	case 0xff:
		if len(b) < 9 {
			return 0, 0
		}

		off := len(b) - 9
		v := binary.BigEndian.Uint64(b[off : off+8])
		return v, 9
	}
}

// PutReverseUint32 encodes a uint32 into the b end and returns the number of bytes written.
func PutReverseUint32(b []byte, v uint32) int {
	switch {
	case v <= 0xfc:
		off := len(b) - 1
		b[off] = byte(v)
		return 1

	case v <= 0xffff:
		off := len(b) - 3
		binary.BigEndian.PutUint16(b[off:off+2], uint16(v))
		b[off+2] = 0xfd
		return 3

	default:
		off := len(b) - 5
		binary.BigEndian.PutUint32(b[off:off+4], uint32(v))
		b[off+4] = 0xfe
		return 5
	}
}

// PutReverseUint64 encodes a uint64 into the b end and returns the number of bytes written.
func PutReverseUint64(b []byte, v uint64) int {
	switch {
	case v <= 0xfc:
		off := len(b) - 1
		b[off] = byte(v)
		return 1

	case v <= 0xffff:
		off := len(b) - 3
		binary.BigEndian.PutUint16(b[off:off+2], uint16(v))
		b[off+2] = 0xfd
		return 3

	case v <= 0xffffffff:
		off := len(b) - 5
		binary.BigEndian.PutUint32(b[off:off+4], uint32(v))
		b[off+4] = 0xfe
		return 5

	default:
		off := len(b) - 9
		binary.BigEndian.PutUint64(b[off:off+8], v)
		b[off+8] = 0xff
		return 9
	}
}
