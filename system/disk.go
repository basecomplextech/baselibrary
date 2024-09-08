// Copyright 2024 Ivan Korobkov. All rights reserved.

package system

import "github.com/basecomplextech/baselibrary/units"

// DiskInfo is a disk information struct.
type DiskInfo struct {
	Total units.Bytes // Total size
	Free  units.Bytes // Free size
	Used  units.Bytes // Used size
}
