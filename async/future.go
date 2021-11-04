package async

// Future returns an async result in the future.
type Future interface {
	// Err returns the future error or nil.
	Err() error

	// Done awaits the completion.
	Done() <-chan struct{}

	// Result returns the current status, result and error.
	Result() (Status, interface{}, error)

	// Cancel tries to cancel the future.
	Cancel() bool
}
