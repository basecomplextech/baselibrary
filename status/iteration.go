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
