package firehose

import (
	"io"
	"sync"
	"time"
)

// Supplier creates a new destination writer (e.g., a file, an upload stream, etc.).
//
// RecordWriter calls Supplier when it needs to start a new "batch" (first write, or
// after a flush/rotation). The returned io.WriteCloser will be closed when a flush
// threshold is reached or when RecordWriter.Close is called.
type Supplier func() (io.WriteCloser, error)

// RecordWriter is a batched writer that "rotates" the underlying io.WriteCloser
// when a threshold is reached.
//
// A flush (rotation) occurs when either:
//   - maxBytes is reached/exceeded (based on bytes successfully written), or
//   - duration has elapsed since the current writer was initialized.
//
// If appendNewLine is true, PutRecord appends '\n' to each record before writing.
//
// Concurrency: all methods are safe for concurrent use; writes and flushes are
// serialized with an internal mutex.
type RecordWriter struct {
	duration       time.Duration
	maxBytes       int
	supplier       Supplier
	appendNewLine  bool
	lock           sync.Mutex
	lastWriter     io.WriteCloser
	estimatedDelay time.Time
	writtenBytes   int
	clock          Clock
}

// PutRecord writes a single record to the current underlying writer.
//
// PutRecord may flush/rotate the underlying writer before and/or after writing
// if a threshold is reached. If appendNewLine is enabled, '\n' is appended to r
// before writing.
func (fw *RecordWriter) PutRecord(r []byte) error {
	fw.lock.Lock()
	defer fw.lock.Unlock()
	// Check for time based threshold before writing
	if err := fw.flushIfThresholdReached(); err != nil {
		return err
	}
	// Initialize writer if not already done
	if err := fw.init(); err != nil {
		return err
	}

	if fw.appendNewLine {
		r = append(r, '\n')
	}
	n, err := fw.lastWriter.Write(r)
	if err != nil {
		return err
	}
	fw.writtenBytes += n

	// Check for max bytes threshold before writing
	if err := fw.flushIfThresholdReached(); err != nil {
		return err
	}
	return nil
}

// Close flushes (closes) the current underlying writer, if any.
//
// After Close, a subsequent PutRecord will lazily create a new writer via Supplier.
func (fw *RecordWriter) Close() error {
	fw.lock.Lock()
	defer fw.lock.Unlock()
	return fw.flush()
}

// FlushIfThresholdReached checks thresholds and flushes (closes) the current
// underlying writer if needed.
//
// This is useful to force time-based rotation in periods without writes (i.e.,
// when PutRecord is not being called).
func (fw *RecordWriter) FlushIfThresholdReached() error {
	fw.lock.Lock()
	defer fw.lock.Unlock()
	return fw.flushIfThresholdReached()
}

func (fw *RecordWriter) init() error {
	if fw.lastWriter == nil {
		writer, err := fw.supplier()
		if err != nil {
			return err
		}
		fw.lastWriter = writer
		fw.writtenBytes = 0
		fw.estimatedDelay = fw.clock.Now().Add(fw.duration)
	}
	return nil
}

func (fw *RecordWriter) flushIfThresholdReached() error {
	if fw.lastWriter != nil && (fw.maxBytes <= fw.writtenBytes || fw.clock.Now().After(fw.estimatedDelay)) {
		return fw.flush()
	}
	return nil
}

func (fw *RecordWriter) flush() error {
	if fw.lastWriter != nil {
		err := fw.lastWriter.Close()
		fw.lastWriter = nil
		return err
	}
	return nil
}
