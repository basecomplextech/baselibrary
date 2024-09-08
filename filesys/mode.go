// Copyright 2022 Ivan Korobkov. All rights reserved.

package filesys

import "os"

// FileMode represents a file's mode and permission bits.
type FileMode = os.FileMode

// The defined file mode bits are the most significant bits of the FileMode.
const (
	// The single letters are the abbreviations
	// Used by the String method's formatting.
	ModeDir        = os.ModeDir        // d: is a directory
	ModeAppend     = os.ModeAppend     // a: append-only
	ModeExclusive  = os.ModeExclusive  // l: exclusive use
	ModeTemporary  = os.ModeTemporary  // T: temporary file; Plan 9 only
	ModeSymlink    = os.ModeSymlink    // L: symbolic link
	ModeDevice     = os.ModeDevice     // D: device file
	ModeNamedPipe  = os.ModeNamedPipe  // p: named pipe (FIFO)
	ModeSocket     = os.ModeSocket     // S: Unix domain socket
	ModeSetuid     = os.ModeSetuid     // u: setuid
	ModeSetgid     = os.ModeSetgid     // g: setgid
	ModeCharDevice = os.ModeCharDevice // c: Unix character device, when ModeDevice is set
	ModeSticky     = os.ModeSticky     // t: sticky
	ModeIrregular  = os.ModeIrregular  // ?: non-regular file; nothing else is known about this file

	// Mask for the type bits. For regular files, none will be set.
	ModeType = os.ModeType

	ModePerm = os.ModePerm // Unix permission bits, 0o777
)
