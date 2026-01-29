package firehose

import (
	"io"
	"sync"
	"time"
)

type Supplier func() (io.WriteCloser, error)

type Writer struct {
	duration       time.Duration
	maxBytes       int
	supplier       Supplier
	appendNewLine  bool
	lock           sync.Mutex
	lastWriter     io.WriteCloser
	estimatedDelay time.Time
	writtenBytes   int
}

func (fw *Writer) PutRecord(r []byte) error {
	fw.lock.Lock()
	defer fw.lock.Unlock()
	// Check for time based threshold before writing
	if err := fw.flushIfThresholdReached(); err != nil {
		return err
	}
	// Initialize writer if not already done
	if err := fw.init(); err != nil {
		return err
	}

	if fw.appendNewLine {
		r = append(r, '\n')
	}
	n, err := fw.lastWriter.Write(r)
	if err != nil {
		return err
	}
	fw.writtenBytes += n

	// Check for max bytes threshold before writing
	if err := fw.flushIfThresholdReached(); err != nil {
		return err
	}
	return nil
}

func (fw *Writer) Close() error {
	fw.lock.Lock()
	defer fw.lock.Unlock()
	return fw.flush()
}

func (fw *Writer) FlushIfThresholdReached() error {
	fw.lock.Lock()
	defer fw.lock.Unlock()
	return fw.flushIfThresholdReached()
}

func (fw *Writer) init() error {
	if fw.lastWriter == nil {
		writer, err := fw.supplier()
		if err != nil {
			return err
		}
		fw.lastWriter = writer
		fw.writtenBytes = 0
		fw.estimatedDelay = time.Now().Add(fw.duration)
	}
	return nil
}

func (fw *Writer) flushIfThresholdReached() error {
	if fw.lastWriter != nil && (fw.maxBytes <= fw.writtenBytes || time.Now().After(fw.estimatedDelay)) {
		return fw.flush()
	}
	return nil
}

func (fw *Writer) flush() error {
	if fw.lastWriter != nil {
		err := fw.lastWriter.Close()
		fw.lastWriter = nil
		return err
	}
	return nil
}
