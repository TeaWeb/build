package teamemory

import (
	"bytes"
	"compress/gzip"
	"github.com/dchest/siphash"
	"github.com/iwind/TeaGo/logs"
	"sync/atomic"
	"unsafe"
)

type ItemType = int8

const (
	ItemInt64     = 1
	ItemBytes     = 2
	ItemInterface = 3
)

func HashKey(key []byte) uint64 {
	return siphash.Hash(0, 0, key)
}

type Item struct {
	Key            []byte
	ExpireAt       int64
	Type           ItemType
	ValueInt64     int64
	ValueBytes     []byte
	ValueInterface interface{}
	IsCompressed   bool

	// linked list
	Prev *Item
	Next *Item

	size int64
}

func NewItem(key []byte, dataType ItemType) *Item {
	return &Item{
		Key:  key,
		Type: dataType,
	}
}

func (this *Item) HashKey() uint64 {
	return HashKey(this.Key)
}

func (this *Item) IncreaseInt64(delta int64) {
	atomic.AddInt64(&this.ValueInt64, delta)
}

func (this *Item) Bytes() []byte {
	if this.IsCompressed {
		reader, err := gzip.NewReader(bytes.NewBuffer(this.ValueBytes))
		if err != nil {
			logs.Error(err)
			return this.ValueBytes
		}

		buf := make([]byte, 256)
		dataBuf := bytes.NewBuffer([]byte{})
		for {
			n, err := reader.Read(buf)
			if n > 0 {
				dataBuf.Write(buf[:n])
			}
			if err != nil {
				break
			}
		}
		return dataBuf.Bytes()
	}
	return this.ValueBytes
}

func (this *Item) String() string {
	return string(this.Bytes())
}

func (this *Item) Size() int64 {
	if this.size == 0 {
		this.size = int64(int(unsafe.Sizeof(*this)) + len(this.Key) + len(this.ValueBytes))
	}
	return this.size
}
