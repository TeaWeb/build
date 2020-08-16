package teacache

import (
	"github.com/TeaWeb/build/internal/teatesting"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestPurgeCache(t *testing.T) {
	if !teatesting.RequireHTTPServer() {
		return
	}

	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:9991/webhook", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Tea-Cache-Purge", "1")
	req.Header.Set("Tea-Key", "z8O4MuXixbKH6aiVyZigYTxxovRblR3u")
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(data))
}
