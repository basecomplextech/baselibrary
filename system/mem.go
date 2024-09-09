// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package system

import "github.com/basecomplextech/baselibrary/units"

type MemoryInfo struct {
	Total units.Bytes
	Free  units.Bytes
	Used  units.Bytes
}
