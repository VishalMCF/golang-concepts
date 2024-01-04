package wal_clone

import (
	"errors"
	"github.com/tidwall/tinylru"
	"os"
	"path/filepath"
	"sync"
)

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

// Options for log
type Options struct {

	// if this is true fsync will not happen. What it means is that the written log file/segment might b stored
	// in the os cache. If this is set to true os will ot flush the changes to the disk and we may lose data if the
	// system goes down
	NoSync bool

	// the log is divided into multiple files. Each file is called a segment. The default segment size is 20MB
	// please search on chatgpt -> why do we need to break the log file into multiple segment
	SegmentSize int

	// whether to write the data into binary format or json format. json format is readable but the file is typically
	// larger as compared to the binary counterpart
	format Logformat

	// segments are also cached in memory to avoid disk look ups. By default at max two segment will be cached
	SegmentCacheSize int

	// ??-:-> stills needs to be check why it is being used
	NoCopy bool

	// Perms represents the datafiles modes and permission bits
	DirPerms  os.FileMode
	FilePerms os.FileMode
}

// DefaultOptions for Open().
var DefaultOptions = &Options{
	NoSync:           false,    // Fsync after every write
	SegmentSize:      20971520, // 20 MB log segment files.
	format:           Binary,   // Binary format is small and fast.
	SegmentCacheSize: 2,        // Number of cached in-memory segments
	NoCopy:           false,    // Make a new copy of data for every Read call.
	DirPerms:         0750,     // Permissions for the created directories
	FilePerms:        0640,     // Permissions for the created data files
}

// Log represents a write ahead log
type Log struct {
	mu         sync.RWMutex
	path       string
	opts       Options
	firstIndex int
	lastIndex  int
	segments   []*segment
	closed     bool
	corrupt    bool
	sfile      *os.File
	wbatch     Batch
	scache     tinylru.LRU
}

// segment represents a single segment file
type segment struct {
	path  string
	index uint64
	ebuf  []byte
	epos  []bpos
}

type bpos struct {
	pos int // byte position
	end int // one byte pos
}

// Open a new write ahead log
func Open(path string, opts *Options) (*Log, error) {
	if opts == nil {
		opts = DefaultOptions
	}
	if opts.SegmentCacheSize <= 0 {
		opts.SegmentCacheSize = DefaultOptions.SegmentCacheSize
	}
	if opts.SegmentSize <= 0 {
		opts.SegmentSize = DefaultOptions.SegmentSize
	}
	if opts.DirPerms == 0 {
		opts.DirPerms = DefaultOptions.DirPerms
	}
	if opts.FilePerms == 0 {
		opts.FilePerms = DefaultOptions.FilePerms
	}

	var err error
	path, err = absPath(path)
	if err != nil {
		return nil, err
	}
	l := &Log{path: path, opts: *opts}
	l.scache.Resize(l.opts.SegmentCacheSize)
	if err := os.MkdirAll(path, l.opts.DirPerms); err != nil {
		return nil, err
	}
	if err := l.load(); err != nil {
		return nil, err
	}
	return l, nil
}

func absPath(path string) (string, error) {
	if path == ":memory:" {
		return "", errors.New("in memory log is not supported")
	}
	return filepath.Abs(path)
}

// Sync performs the fsync on the log. What it means is that wheneber a log is written it will send it to the os to
// write and performs flush over the os cache to persist to the disk. It slowes down the writing operation
func (l *Log) Sync() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.corrupt {
		return ErrCorrupt
	}
	if l.closed {
		return ErrClosed
	}
	return l.sfile.Sync()
}
