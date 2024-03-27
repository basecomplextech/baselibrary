//go:build darwin || dragonfly
// +build darwin dragonfly

package disk

import (
	"syscall"

	"github.com/basecomplextech/baselibrary/units"
)

// GetInfo returns a disk usage info of a directory, e.g. `/`.
func GetInfo(path string) (Info, error) {
	s := syscall.Statfs_t{}
	err := syscall.Statfs(path, &s)
	if err != nil {
		return Info{}, err
	}

	reserved := s.Bfree - s.Bavail // Reserved blocks
	total := uint64(s.Bsize) * (s.Blocks - reserved)
	free := uint64(s.Bsize) * s.Bavail
	used := total - free

	info := Info{
		Total: units.Bytes(total),
		Free:  units.Bytes(free),
		Used:  units.Bytes(used),
	}
	return info, nil
}
