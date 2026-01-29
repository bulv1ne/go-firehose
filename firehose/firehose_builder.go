package firehose

import "time"

type WriterOpts func(*Writer)

func NewFirehoseWriter(supplier Supplier, opts ...WriterOpts) *Writer {
	fw := &Writer{
		duration:      time.Minute,
		maxBytes:      1024 * 1024,
		supplier:      supplier,
		appendNewLine: false,
	}
	for _, opt := range opts {
		opt(fw)
	}
	return fw
}

func WithDuration(d time.Duration) WriterOpts {
	return func(fw *Writer) {
		fw.duration = d
	}
}

func WithMaxBytes(n int) WriterOpts {
	return func(fw *Writer) {
		fw.maxBytes = n
	}
}

func WithAppendNewLine(appendNewLine bool) WriterOpts {
	return func(fw *Writer) {
		fw.appendNewLine = appendNewLine
	}
}
