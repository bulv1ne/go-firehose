package sink

import "bytes"

func NewMemoryByteSink() ByteSink {
	return &MemoryByteSink{}
}

type MemoryByteSink struct {
	buffer bytes.Buffer
}

func (m *MemoryByteSink) Write(p []byte) (n int, err error) {
	return m.buffer.Write(p)
}

func (m *MemoryByteSink) Close() error {
	// No resources to release
	return nil
}

func (m *MemoryByteSink) Bytes() ([]byte, error) {
	return bytes.Clone(m.buffer.Bytes()), nil
}
