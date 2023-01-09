package logging

// Writer writes log records.
type Writer interface {
	// Enabled returns true if the writer is enabled for the given level.
	Enabled(level Level) bool

	// Write writes a record.
	Write(rec Record) error
}

// null

var _ Writer = (*nullWriter)(nil)

type nullWriter struct{}

func newNullWriter() *nullWriter {
	return &nullWriter{}
}

// Enabled returns true if the writer is enabled for the given level.
func (w *nullWriter) Enabled(level Level) bool {
	return false
}

// Write writes a record.
func (w *nullWriter) Write(rec Record) error {
	return nil
}
