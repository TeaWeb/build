package teaproxy

import (
	"bytes"
	"net/http"
	"testing"
)

func TestResponseWriterHeader(t *testing.T) {
	resp := &http.Response{}
	resp.Header = http.Header{
		"Content-Type":  []string{"text/html"},
		"Cache-Control": []string{"no-cache"},
		"Connection":    []string{"Close"},
	}
	resp.ContentLength = 1
	resp.ProtoMajor = 1
	resp.ProtoMinor = 1
	resp.StatusCode = 200

	//resp.Body = ioutil.NopCloser(bytes.NewReader([]byte("hello, world")))

	writer := bytes.NewBuffer([]byte{})
	_ = resp.Write(writer)
	t.Log(string(writer.Bytes()))
}
