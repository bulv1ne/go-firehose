package sink

import (
	"io"
	"os"
)

// NewFileByteSink creates a ByteSink backed by a temporary file.
//
// Closing the returned sink closes and removes the temp file.
func NewFileByteSink() (ByteSink, error) {
	f, err := os.CreateTemp("", "filebytesink_*.tmp")
	if err != nil {
		return nil, err
	}
	return &FileByteSink{file: f}, nil
}

// FileByteSink is a ByteSink implementation that writes to an *os.File and reads
// back from it on demand.
//
// This is useful to avoid holding large payloads in memory.
type FileByteSink struct {
	file *os.File
}

// Write writes bytes to the underlying temp file.
func (f *FileByteSink) Write(p []byte) (n int, err error) {
	return f.file.Write(p)
}

// Close closes the underlying file and removes it from disk.
func (f *FileByteSink) Close() error {
	err := f.file.Close()
	if err != nil {
		return err
	}
	return os.Remove(f.file.Name())
}

// Bytes returns the contents written to the underlying file so far.
//
// It preserves the current file offset by saving it, rewinding to the start to
// read all bytes, and then restoring the offset.
func (f *FileByteSink) Bytes() ([]byte, error) {
	// Save current offset, rewind, read all, then restore offset.
	cur, err := f.file.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}
	if _, err := f.file.Seek(0, io.SeekStart); err != nil {
		return nil, err
	}

	data, err := io.ReadAll(f.file)
	_, _ = f.file.Seek(cur, io.SeekStart)
	return data, err
}
