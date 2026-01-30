# go-firehose

`go-firehose` is a small Go library for **batching records** into an `io.WriteCloser` and **rotating** (closing and recreating) that writer based on:

- **max bytes written** to the current batch, and/or
- **time elapsed** since the current batch started.

It’s useful for log batching, periodic file rotation, chunked uploads, or “write a bunch of records then close/flush” workflows.

## Install

This repo’s module name is `go-firehose` (see `go.mod`), so within this workspace you can import it as:

- `go-firehose/firehose`
- `go-firehose/sink`

If you publish it under a different module path later (e.g. `github.com/<you>/go-firehose`), update `module` in `go.mod` and these imports will follow.

## Core API

### `firehose.Supplier`

A `Supplier` creates a fresh destination writer for a new batch/rotation:

- Called lazily on the first `PutRecord`.
- Called again after a rotation (time/size) or after `Close()`.

```go
type Supplier func() (io.WriteCloser, error)
```

### `firehose.RecordWriter`

A `RecordWriter` writes records to the current underlying writer, and rotates that writer when thresholds are reached.

Key methods:

- `PutRecord([]byte) error` – write one record
- `FlushIfThresholdReached() error` – rotate without writing (useful on a ticker)
- `Close() error` – close the current writer (if any)

Options via `NewRecordWriter(..., opts...)`:

- `WithDuration(d time.Duration)` (default: 1 minute)
- `WithMaxBytes(n int)` (default: 1 MiB)
- `WithAppendNewLine(true|false)` (default: false)

## Quick start

This example batches records into in-memory buffers (see `internal/` in tests for a minimal in-memory `io.WriteCloser` implementation).

```go
package main

import (
	"fmt"
	"io"

	"github.com/bulv1ne/go-firehose/firehose"
)

// Replace this with something real:
// - os.Create(...) to write to real files
// - an upload stream
// - a sink.ByteSink, etc.
func supplier() (io.WriteCloser, error) {
	// TODO: return your io.WriteCloser
	panic("implement me")
}

func main() {
	fw := firehose.NewRecordWriter(
		supplier,
		firehose.WithMaxBytes(1024*1024),
		firehose.WithAppendNewLine(true),
	)
	defer fw.Close()

	_ = fw.PutRecord([]byte("hello"))
	_ = fw.PutRecord([]byte("world"))

	fmt.Println("done")
}
```

## Time-based rotation during idle periods

Time-based rotation is checked during `PutRecord(...)`, but if your system can be idle you may also want a ticker that calls `FlushIfThresholdReached()`.

```go
ticker := time.NewTicker(5 * time.Second)

defer ticker.Stop()

go func() {
	for range ticker.C {
		_ = fw.FlushIfThresholdReached()
	}
}()
```

## Composing writers (gzip, etc.) with `StackWriteCloser`

Many wrappers (like `gzip.Writer`) need to be **closed before** the underlying stream so they can flush trailers/footers.

`firehose.StackWriteCloser` helps track close order:

- `Push(...)` sets the current write target.
- `Close()` closes everything in **reverse (LIFO)** order.

Example supplier that gzips each batch:

```go
func gzipSupplier() (io.WriteCloser, error) {
	stack := firehose.NewWriteCloserStack()

	// Push the underlying destination first.
	// (This is the writer that ultimately receives bytes.)
	base := stack.Push(/* your io.WriteCloser */)

	// Push the wrapper second. gzip.NewWriter takes io.Writer and returns *gzip.Writer
	// which implements io.WriteCloser.
	stack.Push(gzip.NewWriter(base))

	return stack, nil
}
```

## Package layout

- `firehose/` – main library (`RecordWriter`, `Supplier`, writer stack)
- `sink/` – helper sink abstractions/implementations
- `internal/` – test helpers (not meant for external use)

## Development

Run tests:

```sh
go test ./...
```
