package firehose

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryWriteCloser(t *testing.T) {
	resetHashMap()

	mwc := NewMemoryWriteCloser()
	data := []byte("Hello, Firehose!")
	_, err := mwc.Write(data)

	assert.NoError(t, err)
	assert.NoError(t, mwc.Close())

	assert.Equal(t, int32(1), nameCounter)
	assert.Equal(t, mwcHashMap["1"], data)
}
