package async

// Status is a future status.
type Status int

const (
	StatusPending Status = iota
	StatusOK             // completed
	StatusError          // failed with an error
	StatusExit           // cancelled or exited without result
)

func (s Status) String() string {
	switch s {
	case StatusPending:
		return "pending"
	case StatusOK:
		return "ok"
	case StatusError:
		return "error"
	case StatusExit:
		return "exit"
	}
	return ""
}
