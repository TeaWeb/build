package configs

import (
	"bytes"
	"fmt"
	"github.com/iwind/TeaGo/types"
	"io/ioutil"
)

// id|sum|flags|data
type Item struct {
	Id    string
	Sum   string
	Flags []int
	Data  []byte
}

func UnmarshalItem(data []byte) *Item {
	item := &Item{}

	buf := bytes.NewBuffer(data)

	{
		line, err := buf.ReadBytes('|')
		if err != nil {
			return item
		}
		item.Id = string(line[:len(line)-1])
	}

	{
		line, err := buf.ReadBytes('|')
		if err != nil {
			return item
		}
		item.Sum = string(line[:len(line)-1])
	}

	{
		line, err := buf.ReadBytes('|')
		if err != nil {
			return item
		}
		line = line[:len(line)-1]
		if len(line) > 0 {
			pieces := bytes.Split(line, []byte{','})
			for _, piece := range pieces {
				if len(piece) > 0 {
					item.Flags = append(item.Flags, types.Int(piece))
				}
			}
		}
	}

	result, err := ioutil.ReadAll(buf)
	if err == nil {
		item.Data = result
	}

	return item
}

func NewItem() *Item {
	return &Item{}
}

func (this *Item) Marshal() []byte {
	data := []byte(this.Id)
	data = append(data, '|')
	data = append(data, this.Sum...)
	data = append(data, '|')
	for _, flag := range this.Flags {
		data = append(data, fmt.Sprintf("%d,", flag)...)
	}
	data = append(data, '|')
	data = append(data, this.Data...)
	return data
}
