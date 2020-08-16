package teacache

import "encoding/binary"

// 缓存条目
type Item struct {
	Header []byte
	Body   []byte
}

// 获取新对象
func NewItem() *Item {
	return &Item{}
}

func (this *Item) Encode() (data []byte) {
	l := make([]byte, 8)
	binary.BigEndian.PutUint32(l, uint32(len(this.Header)))
	l = append(l, this.Header ...)
	l = append(l, this.Body ...)
	return l
}

func (this *Item) Decode(data []byte) {
	l := data[:8]
	headerLength := binary.BigEndian.Uint32(l)
	this.Header = data[8 : 8+headerLength]
	this.Body = data[8+headerLength:]
}

