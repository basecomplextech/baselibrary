package logging2

// Writer writes log records.
type Writer interface {
	// Write writes a record.
	Write(record Record) error
}
