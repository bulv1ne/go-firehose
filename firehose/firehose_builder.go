package firehose

import "time"

type WriterOpts func(*RecordWriter)

func NewRecordWriter(supplier Supplier, opts ...WriterOpts) *RecordWriter {
	fw := &RecordWriter{
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
	return func(fw *RecordWriter) {
		fw.duration = d
	}
}

func WithMaxBytes(n int) WriterOpts {
	return func(fw *RecordWriter) {
		fw.maxBytes = n
	}
}

func WithAppendNewLine(appendNewLine bool) WriterOpts {
	return func(fw *RecordWriter) {
		fw.appendNewLine = appendNewLine
	}
}
