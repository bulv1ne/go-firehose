package sink

import "bytes"

// NewMemoryByteSink creates an in-memory ByteSink.
//
// The returned sink buffers all written data in memory.
func NewMemoryByteSink() ByteSink {
	return &MemoryByteSink{}
}

// MemoryByteSink is a ByteSink implementation that stores all bytes in a bytes.Buffer.
type MemoryByteSink struct {
	buffer bytes.Buffer
}

// Write appends p to the in-memory buffer.
func (m *MemoryByteSink) Write(p []byte) (n int, err error) {
	return m.buffer.Write(p)
}

// Close releases no resources and always returns nil.
func (m *MemoryByteSink) Close() error {
	// No resources to release
	return nil
}

// Bytes returns a copy of the buffered bytes.
func (m *MemoryByteSink) Bytes() ([]byte, error) {
	return bytes.Clone(m.buffer.Bytes()), nil
}
