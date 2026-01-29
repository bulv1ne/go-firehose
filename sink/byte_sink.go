package sink

import "io"

type ByteSink interface {
	io.WriteCloser
	Bytes() ([]byte, error)
}
