package filesys

import "os"

var (
	ErrInvalid    = os.ErrInvalid
	ErrPermission = os.ErrPermission
	ErrExist      = os.ErrExist
	ErrNotExist   = os.ErrNotExist
	ErrClosed     = os.ErrClosed
)
