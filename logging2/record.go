package logging2

type Record struct {
	Level   Level
	Logger  string
	Message string
	Fields  []Field
	Stack   []byte
}

// NewRecord returns a new record with info level.
func NewRecord(logger string) Record {
	return Record{
		Level:  LevelInfo,
		Logger: logger,
	}
}
