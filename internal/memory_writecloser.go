package internal

import (
	"fmt"
	"io"
	"sync"
	"sync/atomic"

	"github.com/bulv1ne/go-firehose/sink"
)

var nameCounter int32 = 0
var mwcHashMap = make(map[string][]byte)
var mwcHashMapLock sync.Mutex

func ResetHashMap() {
	atomic.StoreInt32(&nameCounter, 0)
	mwcHashMapLock.Lock()
	mwcHashMap = make(map[string][]byte)
	mwcHashMapLock.Unlock()
}

func GetHashMap() map[string][]byte {
	mwcHashMapLock.Lock()
	defer mwcHashMapLock.Unlock()
	// Return a copy to avoid
	copyMap := make(map[string][]byte)
	for k, v := range mwcHashMap {
		copyMap[k] = v
	}
	return copyMap
}

type MemoryWriteCloser struct {
	name string
	sink sink.ByteSink
}

func NewMemoryWriteCloser() io.WriteCloser {
	atomic.AddInt32(&nameCounter, 1)
	return &MemoryWriteCloser{
		name: fmt.Sprintf("%x", nameCounter),
		sink: sink.NewMemoryByteSink(),
	}
}

func (mwc *MemoryWriteCloser) Write(p []byte) (n int, err error) {
	return mwc.sink.Write(p)
}

func (mwc *MemoryWriteCloser) Close() error {
	bytes, err := mwc.sink.Bytes()
	if err != nil {
		return err
	}
	mwcHashMapLock.Lock()
	mwcHashMap[mwc.name] = bytes
	mwcHashMapLock.Unlock()
	return mwc.sink.Close()
}
