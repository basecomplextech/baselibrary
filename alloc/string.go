// Copyright 2024 Ivan Korobkov. All rights reserved.

package alloc

import (
	"fmt"
	"unicode/utf8"
	"unsafe"

	"github.com/basecomplextech/baselibrary/alloc/internal/heap"
)

// String allocates a new string and copies data from src into it.
func String(a Arena, src string) string {
	if len(src) == 0 {
		return ""
	}

	dst := Bytes(a, len(src))
	copy(dst, src)
	return unsafeString(dst)
}

// StringBytes allocates a new string and copies data from src into it.
func StringBytes(a Arena, src []byte) string {
	if len(src) == 0 {
		return ""
	}

	dst := Bytes(a, len(src))
	copy(dst, src)
	return unsafeString(dst)
}

// StringRunes allocates a new string and copies data from src into it.
func StringRunes(a Arena, src []rune) string {
	size := 0
	for _, r := range src {
		size += utf8.RuneLen(r)
	}

	dst := allocSlice[[]byte](a, size, size)
	n := 0

	for _, r := range src {
		n += utf8.EncodeRune(dst[n:], r)
	}

	return unsafeString(dst)
}

// StringJoin joins strings using separator.
func StringJoin(a Arena, src []string, sep string) string {
	size := 0
	for i, s := range src {
		if i > 0 {
			size += len(sep)
		}
		size += len(s)
	}
	if size == 0 {
		return ""
	}

	b := allocSlice[[]byte](a, 0, size)
	for i, s := range src {
		if i > 0 {
			b = append(b, sep...)
		}
		b = append(b, s...)
	}

	return unsafeString(b)
}

// StringJoin2 joins two strings using separator.
func StringJoin2(a Arena, s1, s2, sep string) string {
	size := len(s1) + len(s2) + len(sep)
	if size == 0 {
		return ""
	}

	b := allocSlice[[]byte](a, 0, size)
	b = append(b, s1...)
	b = append(b, sep...)
	b = append(b, s2...)

	return unsafeString(b)
}

// StringFormat formats a string using fmt.Appendf and returns a new string allocated in the arena.
func StringFormat(a Arena, format string, args ...any) string {
	// Allocate scratch block
	size := heap.MinBlockSize
	if len(format) > size {
		size = len(format)
	}

	block := heap.Global.Alloc(size)
	defer heap.Global.Free(block)

	// Format string
	b := block.Bytes()
	b = fmt.Appendf(b, format, args...)

	// Copy result
	return StringBytes(a, b)
}

// private

func unsafeString(b []byte) string {
	if b == nil {
		return ""
	}
	return *(*string)(unsafe.Pointer(&b))
}
