package disk

import "github.com/basecomplextech/baselibrary/units"

// Info is a disk stats struct.
type Info struct {
	Total units.Bytes // Total size
	Free  units.Bytes // Free size
	Used  units.Bytes // Used size
}
