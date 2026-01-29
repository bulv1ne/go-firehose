package firehose

import (
	"errors"
	"io"
)

// NewWriteCloserStack constructs an empty StackWriteCloser.
//
// Call Push at least once before calling Write, otherwise Write will panic
// due to a nil current writer. Close is safe to call even if nothing was pushed.
func NewWriteCloserStack() *StackWriteCloser {
	return &StackWriteCloser{}
}

// StackWriteCloser is a minimal "stack" of io.WriteCloser values.
//
// The most recently pushed writer becomes the current write target. Close closes
// all pushed writers in reverse (LIFO) order and joins any close errors.
//
// Note: This type does not automatically wrap writers; it only tracks the order
// you push them and routes writes to the latest one.
type StackWriteCloser struct {
	// stack holds all pushed write-closers in push order.
	stack []io.WriteCloser
	// current is the write target used by Write; it is set by Push.
	current io.WriteCloser
}

// Write forwards p to the most recently pushed io.WriteCloser.
//
// Push must have been called at least once before Write is used.
func (wcs *StackWriteCloser) Write(p []byte) (n int, err error) {
	return wcs.current.Write(p)
}

// Close closes all pushed io.WriteCloser values in reverse (LIFO) order.
//
// If multiple closes fail, the returned error is the joined error value.
// If nothing was pushed, Close returns nil.
func (wcs *StackWriteCloser) Close() error {
	errorList := make([]error, 0)
	for i := len(wcs.stack) - 1; i >= 0; i-- {
		if err := wcs.stack[i].Close(); err != nil {
			errorList = append(errorList, err)
		}
	}
	return errors.Join(errorList...)
}

// Push appends wc to the stack and sets it as the current write target.
//
// The returned value is wc for convenience (e.g., inline assignment).
func (wcs *StackWriteCloser) Push(wc io.WriteCloser) io.WriteCloser {
	wcs.stack = append(wcs.stack, wc)
	wcs.current = wc
	return wc
}
