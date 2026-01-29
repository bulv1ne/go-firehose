package sink

import (
	"io"
	"os"
)

func NewFileByteSink() (ByteSink, error) {
	f, err := os.CreateTemp("", "filebytesink_*.tmp")
	if err != nil {
		return nil, err
	}
	return &FileByteSink{file: f}, nil
}

type FileByteSink struct {
	file *os.File
}

func (f *FileByteSink) Write(p []byte) (n int, err error) {
	return f.file.Write(p)
}

func (f *FileByteSink) Close() error {
	err := f.file.Close()
	if err != nil {
		return err
	}
	return os.Remove(f.file.Name())
}

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
