package compactint

import (
	"encoding/binary"
)

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

// PutReverseUint32 appends a uint32 to b and returns the number of bytes written.
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

// PutReverseUint64 appends a uint64 to b and returns the number of bytes written.
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
