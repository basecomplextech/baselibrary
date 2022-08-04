package logging

import "time"

type Record struct {
	Time    time.Time `json:"time"`
	Level   Level     `json:"level"`
	Logger  string    `json:"logger"`
	Message string    `json:"message"`
	Fields  []Field   `json:"fields"`
	Stack   []byte    `json:"stack"`
}

// NewRecord returns a new record with info level.
func NewRecord(logger string) Record {
	return Record{
		Time:   time.Now(),
		Level:  LevelInfo,
		Logger: logger,
	}
}
