// Copyright 2024 Ivan Korobkov. All rights reserved.

//go:build darwin || dragonfly
// +build darwin dragonfly

package system

import (
	"syscall"

	"github.com/basecomplextech/baselibrary/units"
)

// Disk returns a disk usage info of a directory, e.g. `/`.
func Disk(path string) (DiskInfo, error) {
	s := syscall.Statfs_t{}
	err := syscall.Statfs(path, &s)
	if err != nil {
		return DiskInfo{}, err
	}

	reserved := s.Bfree - s.Bavail // Reserved blocks
	total := uint64(s.Bsize) * (s.Blocks - reserved)
	free := uint64(s.Bsize) * s.Bavail
	used := total - free

	info := DiskInfo{
		Total: units.Bytes(total),
		Free:  units.Bytes(free),
		Used:  units.Bytes(used),
	}
	return info, nil
}
