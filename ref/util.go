// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package ref

import (
	_ "unsafe"
)

//go:linkname fastrand runtime.fastrand
func fastrand() uint32
