package system

import "github.com/basecomplextech/baselibrary/units"

type MemoryInfo struct {
	Total units.Bytes
	Free  units.Bytes
	Used  units.Bytes
}
