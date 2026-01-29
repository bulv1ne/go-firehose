package firehose

import (
	"bytes"
	"compress/gzip"
	"io"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func MemoryWriteCloserSupplier() (io.WriteCloser, error) {
	return NewMemoryWriteCloser(), nil
}

func GzipMemoryWriteCloserSupplier() (io.WriteCloser, error) {
	stack := NewWriteCloserStack()
	writer := stack.Push(NewMemoryWriteCloser())
	stack.Push(gzip.NewWriter(writer))
	return stack, nil
}

func MapKeys(m map[string][]byte) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func TestFirehoseWriter_Delay(t *testing.T) {
	resetHashMap()

	fw := NewFirehoseWriter(MemoryWriteCloserSupplier, WithDuration(20*time.Millisecond))

	assert.NoError(t, fw.PutRecord([]byte("Niels")))
	assert.NoError(t, fw.PutRecord([]byte("Tisse")))

	time.Sleep(30 * time.Millisecond)
	assert.NoError(t, fw.FlushIfThresholdReached())

	assert.NoError(t, fw.PutRecord([]byte("Juna")))
	assert.NoError(t, fw.PutRecord([]byte("Alise")))

	assert.NoError(t, fw.Close())

	assert.Equal(t, MapKeys(mwcHashMap), []string{"1", "2"})
	assert.Equal(t, mwcHashMap["1"], []byte("NielsTisse"))
	assert.Equal(t, mwcHashMap["2"], []byte("JunaAlise"))
}

func TestFirehoseWriter_MaxBytes(t *testing.T) {
	resetHashMap()

	fw := NewFirehoseWriter(MemoryWriteCloserSupplier, WithMaxBytes(7))

	assert.NoError(t, fw.PutRecord([]byte("Niels")))
	assert.NoError(t, fw.PutRecord([]byte("Tisse")))
	assert.NoError(t, fw.PutRecord([]byte("Juna")))
	assert.NoError(t, fw.PutRecord([]byte("Alise")))

	assert.NoError(t, fw.Close())

	assert.Equal(t, MapKeys(mwcHashMap), []string{"1", "2"})
	assert.Equal(t, mwcHashMap["1"], []byte("NielsTisse"))
	assert.Equal(t, mwcHashMap["2"], []byte("JunaAlise"))
}

func TestFirehoseWriter_AppendNewLine(t *testing.T) {
	resetHashMap()

	fw := NewFirehoseWriter(MemoryWriteCloserSupplier, WithAppendNewLine(true))

	assert.NoError(t, fw.PutRecord([]byte("Niels")))
	assert.NoError(t, fw.PutRecord([]byte("Tisse")))
	assert.NoError(t, fw.PutRecord([]byte("Juna")))
	assert.NoError(t, fw.PutRecord([]byte("Alise")))

	assert.NoError(t, fw.Close())

	assert.Equal(t, MapKeys(mwcHashMap), []string{"1"})
	assert.Equal(t, mwcHashMap["1"], []byte("Niels\nTisse\nJuna\nAlise\n"))
}

func TestFirehoseWriter_Gzip(t *testing.T) {
	resetHashMap()

	fw := NewFirehoseWriter(GzipMemoryWriteCloserSupplier)

	assert.NoError(t, fw.PutRecord([]byte("Niels")))
	assert.NoError(t, fw.PutRecord([]byte("Tisse")))
	assert.NoError(t, fw.PutRecord([]byte("Juna")))
	assert.NoError(t, fw.PutRecord([]byte("Alise")))

	assert.NoError(t, fw.Close())

	assert.Equal(t, MapKeys(mwcHashMap), []string{"1"})
	gzippedData := mwcHashMap["1"]
	reader, err := gzip.NewReader(bytes.NewReader(gzippedData))
	assert.NoError(t, err)
	data, err := io.ReadAll(reader)
	assert.NoError(t, err)
	assert.Equal(t, data, []byte("NielsTisseJunaAlise"))
}
