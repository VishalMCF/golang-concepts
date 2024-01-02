package main

import "errors"

// defining all the errors possible

var (
	// ErrCorrupt is rturned when the log is corrupt
	ErrCorrupt = errors.New("log corrupt")

	// ErrClosed is returned when an operation cannot be completed because the log is closed
	ErrClosed = errors.New("log is closed ")

	// ErrNotFound is returned when an entry is not found
	ErrNotFound = errors.New("log not found")

	// ErrOutofOrder is returned when the index on which Write() is called is not equal to the LastIndex()+1
	// It is required that the log monotonically grows and does not contain any gaps.
	// This order is valid :- 10,11,12,13,14,15
	// This order is not valid :- 11,12,13,15 && 11,13,12,14,15
	ErrOutofOrder = errors.New("log corrupt")

	// ErrOutOfRange is returned from TruncateFront() and TruncateBack() when
	// the index not in the range of the log's first and last index. Or, this
	// may be returned when the caller is attempting to remove *all* entries;
	// The log requires that at least one entry exists following a truncate.
	ErrOutOfRange = errors.New("out of range")
)

type Logformat byte

const (
	Binary Logformat = 0
	JSON   Logformat = 1
)
