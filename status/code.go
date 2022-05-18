package status

type Code string

const (
	CodeUndefined   Code = ""
	CodeError       Code = "error"
	CodeOK          Code = "ok"
	CodeStopped     Code = "stopped"
	CodeTimeout     Code = "timeout"
	CodeUnavailable Code = "unavailable"
	CodeNotFound    Code = "not_found"
	CodeCorrupted   Code = "corrupted"
	CodeTest        Code = "test"
)
