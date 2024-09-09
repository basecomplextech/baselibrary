// Copyright 2023 Ivan Korobkov. All rights reserved.
// Use of this software is governed by the MIT License
// that can be found in the LICENSE file.

package status

import "fmt"

var (
	End  = New(CodeEnd, "")
	Wait = New(CodeWait, "")
)

func Endf(format string, a ...any) Status {
	msg := fmt.Sprintf(format, a...)
	return Status{Code: CodeEnd, Message: msg}
}
