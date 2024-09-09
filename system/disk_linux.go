// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

//go:build linux && !s390x && !arm && !386
// +build linux,!s390x,!arm,!386

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
	total := uint64(s.Frsize) * (s.Blocks - reserved)
	free := uint64(s.Frsize) * s.Bavail
	used := total - free

	info := DiskInfo{
		Total: units.Bytes(total),
		Free:  units.Bytes(free),
		Used:  units.Bytes(used),
	}
	return info, nil
}
