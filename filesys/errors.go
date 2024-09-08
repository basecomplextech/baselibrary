// Copyright 2022 Ivan Korobkov. All rights reserved.

package filesys

import "os"

var (
	ErrInvalid    = os.ErrInvalid
	ErrPermission = os.ErrPermission
	ErrExist      = os.ErrExist
	ErrNotExist   = os.ErrNotExist
	ErrClosed     = os.ErrClosed
)
