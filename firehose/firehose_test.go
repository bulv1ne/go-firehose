package firehose

import (
	"bytes"
	"compress/gzip"
	"io"
	"sort"
	"testing"
	"time"

	"github.com/bulv1ne/go-firehose/internal"

	"github.com/stretchr/testify/assert"
)

func MemoryWriteCloserSupplier() (io.WriteCloser, error) {
	return internal.NewMemoryWriteCloser(), nil
}

func GzipMemoryWriteCloserSupplier() (io.WriteCloser, error) {
	stack := NewWriteCloserStack()
	writer := stack.Push(internal.NewMemoryWriteCloser())
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

func TestRecordWriter_Delay(t *testing.T) {
	internal.ResetHashMap()

	clock := &MockClock{CurrentTime: time.Now()}
	fw := NewRecordWriter(MemoryWriteCloserSupplier, WithDuration(1*time.Minute), WithClock(clock))

	assert.NoError(t, fw.PutRecord([]byte("Niels")))
	assert.NoError(t, fw.PutRecord([]byte("Tisse")))

	clock.Advance(2 * time.Minute)
	assert.NoError(t, fw.FlushIfThresholdReached())

	assert.NoError(t, fw.PutRecord([]byte("Juna")))
	assert.NoError(t, fw.PutRecord([]byte("Alise")))

	assert.NoError(t, fw.Close())

	hashMap := internal.GetHashMap()
	assert.Equal(t, MapKeys(hashMap), []string{"1", "2"})
	assert.Equal(t, hashMap["1"], []byte("NielsTisse"))
	assert.Equal(t, hashMap["2"], []byte("JunaAlise"))
}

func TestRecordWriter_MaxBytes(t *testing.T) {
	internal.ResetHashMap()

	fw := NewRecordWriter(MemoryWriteCloserSupplier, WithMaxBytes(7))

	assert.NoError(t, fw.PutRecord([]byte("Niels")))
	assert.NoError(t, fw.PutRecord([]byte("Tisse")))
	assert.NoError(t, fw.PutRecord([]byte("Juna")))
	assert.NoError(t, fw.PutRecord([]byte("Alise")))

	assert.NoError(t, fw.Close())

	hashMap := internal.GetHashMap()
	assert.Equal(t, MapKeys(hashMap), []string{"1", "2"})
	assert.Equal(t, hashMap["1"], []byte("NielsTisse"))
	assert.Equal(t, hashMap["2"], []byte("JunaAlise"))
}

func TestRecordWriter_AppendNewLine(t *testing.T) {
	internal.ResetHashMap()

	fw := NewRecordWriter(MemoryWriteCloserSupplier, WithAppendNewLine(true))

	assert.NoError(t, fw.PutRecord([]byte("Niels")))
	assert.NoError(t, fw.PutRecord([]byte("Tisse")))
	assert.NoError(t, fw.PutRecord([]byte("Juna")))
	assert.NoError(t, fw.PutRecord([]byte("Alise")))

	assert.NoError(t, fw.Close())

	hashMap := internal.GetHashMap()
	assert.Equal(t, MapKeys(hashMap), []string{"1"})
	assert.Equal(t, hashMap["1"], []byte("Niels\nTisse\nJuna\nAlise\n"))
}

func TestRecordWriter_Gzip(t *testing.T) {
	internal.ResetHashMap()

	fw := NewRecordWriter(GzipMemoryWriteCloserSupplier)

	assert.NoError(t, fw.PutRecord([]byte("Niels")))
	assert.NoError(t, fw.PutRecord([]byte("Tisse")))
	assert.NoError(t, fw.PutRecord([]byte("Juna")))
	assert.NoError(t, fw.PutRecord([]byte("Alise")))

	assert.NoError(t, fw.Close())

	hashMap := internal.GetHashMap()
	assert.Equal(t, MapKeys(hashMap), []string{"1"})
	gzippedData := hashMap["1"]
	reader, err := gzip.NewReader(bytes.NewReader(gzippedData))
	assert.NoError(t, err)
	data, err := io.ReadAll(reader)
	assert.NoError(t, err)
	assert.Equal(t, data, []byte("NielsTisseJunaAlise"))
}
