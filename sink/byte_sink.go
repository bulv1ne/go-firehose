package sink

import "io"

// ByteSink is a write target that can later return the bytes written to it.
//
// Implementations may buffer in memory or spool to disk. Bytes should return the
// data written so far (subject to the implementation's semantics) and may fail
// if the underlying storage cannot be read.
type ByteSink interface {
	io.WriteCloser
	Bytes() ([]byte, error)
}
