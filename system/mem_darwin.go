// Copyright 2024 Ivan Korobkov. All rights reserved.

//go:build darwin
// +build darwin

package system

import (
	"os/exec"
	"regexp"
	"strconv"
	"syscall"
	"unsafe"

	"github.com/basecomplextech/baselibrary/units"
)

// Memory returns the system memory information.
func Memory() (MemoryInfo, error) {
	total, err := totalMemory()
	if err != nil {
		return MemoryInfo{}, err
	}
	free, err := freeMemory()
	if err != nil {
		return MemoryInfo{}, err
	}
	used := total - free

	info := MemoryInfo{
		Total: units.Bytes(total),
		Free:  units.Bytes(free),
		Used:  units.Bytes(used),
	}
	return info, nil
}

// private

func totalMemory() (int64, error) {
	v, err := sysctlUint64("hw.memsize")
	if err != nil {
		return 0, err
	}
	return int64(v), nil
}

func freeMemory() (int64, error) {
	cmd := exec.Command("vm_stat")
	out, err := cmd.Output()
	if err != nil {
		return 0, nil
	}

	pageSizeRegex := regexp.MustCompile("page size of ([0-9]*) bytes")
	freePagesRegex := regexp.MustCompile("Pages free: *([0-9]*)\\.")

	// Parse page size
	pageSize := int64(0)
	matches := pageSizeRegex.FindSubmatchIndex(out)
	if len(matches) == 4 {
		part := string(out[matches[2]:matches[3]])
		pageSize, err = strconv.ParseInt(part, 10, 64)
		if err != nil {
			return 0, err
		}
	}

	// Parse free pages
	freePages := int64(0)
	matches = freePagesRegex.FindSubmatchIndex(out)
	if len(matches) == 4 {
		part := string(out[matches[2]:matches[3]])
		freePages, err = strconv.ParseInt(part, 10, 64)
		if err != nil {
			return 0, err
		}
	}

	free := freePages * pageSize
	return free, nil
}

func sysctlUint64(name string) (uint64, error) {
	s, err := syscall.Sysctl(name)
	if err != nil {
		return 0, err
	}

	// hack because the string conversion above drops a \0
	b := []byte(s)
	if len(b) < 8 {
		b = append(b, 0)
	}

	return *(*uint64)(unsafe.Pointer(&b[0])), nil
}
