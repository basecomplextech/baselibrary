// Copyright 2022 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package filesys

import "os"

var (
	ErrInvalid    = os.ErrInvalid
	ErrPermission = os.ErrPermission
	ErrExist      = os.ErrExist
	ErrNotExist   = os.ErrNotExist
	ErrClosed     = os.ErrClosed
)
