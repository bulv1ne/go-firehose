package sink

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertByteSink(t *testing.T, sink ByteSink) {
	data1 := []byte("Hello, Memory Byte Sink!")
	data2 := []byte("Bye, File Byte Sink!")
	expected := []byte("Hello, Memory Byte Sink!Bye, File Byte Sink!")

	n, err := sink.Write(data1)
	assert.NoError(t, err)
	assert.Equal(t, len(data1), n)

	n, err = sink.Write(data2)
	assert.NoError(t, err)
	assert.Equal(t, len(data2), n)

	result, err := sink.Bytes()
	assert.NoError(t, err)
	assert.Equal(t, expected, result)

	assert.NoError(t, sink.Close())
}

func TestMemoryByteSink_Bytes(t *testing.T) {
	sink := NewMemoryByteSink()

	AssertByteSink(t, sink)
}

func TestFileByteSink_Bytes(t *testing.T) {
	sink, err := NewFileByteSink()
	assert.NoError(t, err)

	AssertByteSink(t, sink)
}
