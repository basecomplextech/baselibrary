// Copyright 2024 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package status

import (
	"fmt"
)

// ParseError

// ParseError returns a data corruption error status.
func ParseError(msg string) Status {
	return Status{
		Code:    CodeParseError,
		Message: msg,
	}
}

// ParseErrorf formats and returns a data corruption error status.
func ParseErrorf(format string, a ...any) Status {
	msg := format
	if len(a) > 0 {
		msg = fmt.Sprintf(format, a...)
	}

	return Status{
		Code:    CodeParseError,
		Message: msg,
	}
}

// WrapParseError returns a data corruption error status.
func WrapParseError(err error) Status {
	return wrapError(err, CodeParseError)
}

// WrapParseErrorf formats and returns a data corruption error status.
func WrapParseErrorf(err error, format string, a ...any) Status {
	return wrapErrorf(err, CodeParseError, format, a...)
}

// ChecksumError

// ChecksumError returns a data corruption error status.
func ChecksumError(msg string) Status {
	return Status{
		Code:    CodeChecksumError,
		Message: msg,
	}
}

// ChecksumErrorf formats and returns a data corruption error status.
func ChecksumErrorf(format string, a ...any) Status {
	msg := format
	if len(a) > 0 {
		msg = fmt.Sprintf(format, a...)
	}

	return Status{
		Code:    CodeChecksumError,
		Message: msg,
	}
}

// WrapChecksumError returns a data corruption error status.
func WrapChecksumError(err error) Status {
	return wrapError(err, CodeChecksumError)
}

// WrapChecksumErrorf formats and returns a data corruption error status.
func WrapChecksumErrorf(err error, format string, a ...any) Status {
	return wrapErrorf(err, CodeChecksumError, format, a...)
}
