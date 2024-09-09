// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

//go:build linux
// +build linux

package system

import (
	"syscall"

	"github.com/basecomplextech/baselibrary/units"
)

// Memory returns the system memory information.
func Memory() (MemoryInfo, error) {
	info := &syscall.Sysinfo_t{}
	if err := syscall.Sysinfo(info); err != nil {
		return MemoryInfo{}, err
	}

	total := int64(info.Totalram) * int64(info.Unit)
	free := int64(info.Freeram) * int64(info.Unit)
	used := total - free

	result := MemoryInfo{
		Total: units.Bytes(total),
		Free:  units.Bytes(free),
		Used:  units.Bytes(used),
	}
	return result, nil
}
