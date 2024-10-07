// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package status

var _ error = (*Err)(nil)

// Err is a status error.
type Err struct {
	Code    Code
	Message string
	Cause   error
}

// ToError converts a status into an error, or returns nil if the status is OK.
func ToError(s Status) error {
	if s.Code == CodeOK {
		return nil
	}

	return &Err{
		Code:    s.Code,
		Message: s.Message,
		Cause:   s.Error,
	}
}

// Error implements the error interface.
func (e *Err) Error() string {
	return e.Message
}

// Status converts the error into a status.
func (e *Err) Status() Status {
	return Status{
		Code:    e.Code,
		Message: e.Message,
		Error:   e.Cause,
	}
}
