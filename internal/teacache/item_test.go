package teacache

import (
	"bufio"
	"bytes"
	"github.com/iwind/TeaGo/assert"
	"net/http"
	"testing"
)

func TestItem_Encode(t *testing.T) {
	item := &Item{
		Header: []byte("Hello"),
		Body:   []byte("World"),
	}

	a := assert.NewAssertion(t).Quiet()
	a.Log(string(item.Encode()))

	newItem := &Item{}
	newItem.Decode(item.Encode())
	a.Equals(string(newItem.Header), "Hello")
	a.Equals(string(newItem.Body), "World")
}

func TestItem_Response(t *testing.T) {
	item := NewItem()
	item.Header = []byte(`HTTP/1.1 200 OK
Content-Type: image/png
Etag: "et282346d6373bcae13d89ac46447a228c"
Last-Modified: Fri, 19 Apr 2019 08:01:32 GMT

`)
	resp, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(item.Header)), nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp.Header)
}