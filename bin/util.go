// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package bin

import "unsafe"

func unsafeByteString(s string) []byte {
	if s == "" {
		return nil
	}

	d := unsafe.StringData(s)
	return unsafe.Slice(d, len(s))
}

//go:linkname fastrand runtime.fastrand
func fastrand() uint32
