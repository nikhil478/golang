package wal

import "time"

type Operation string

const (
	OPERATION_SET    Operation = "SET"
	OPERATION_DELETE Operation = "DELETE"
)

type Entry struct {
	Sequence  uint64
	Timestamp time.Time
	Operation Operation
	Key       string
	Data      []byte
}

type WAL interface {
	Append(entry Entry) error
	Replay() ([]Entry, error)
}
