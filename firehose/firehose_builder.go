package firehose

import "time"

type WriterOpts func(*RecordWriter)

func NewRecordWriter(supplier Supplier, opts ...WriterOpts) *RecordWriter {
	fw := &RecordWriter{
		duration:      time.Minute,
		maxBytes:      1024 * 1024,
		supplier:      supplier,
		appendNewLine: false,
		clock:         RealClock{},
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

// WithClock sets a custom clock for time-based threshold checks.
//
// Useful for testing.
func WithClock(clock Clock) WriterOpts {
	return func(fw *RecordWriter) {
		fw.clock = clock
	}
}
