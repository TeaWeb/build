package checkpoints

import (
	"bytes"
	"github.com/TeaWeb/build/internal/teawaf/requests"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
)

func TestRequestFormArgCheckpoint_RequestValue(t *testing.T) {
	rawReq, err := http.NewRequest(http.MethodPost, "http://teaos.cn", bytes.NewBuffer([]byte("name=lu&age=20&encoded="+url.QueryEscape("<strong>ENCODED STRING</strong>"))))
	if err != nil {
		t.Fatal(err)
	}

	req := requests.NewRequest(rawReq)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	checkpoint := new(RequestFormArgCheckpoint)
	t.Log(checkpoint.RequestValue(req, "name", nil))
	t.Log(checkpoint.RequestValue(req, "age", nil))
	t.Log(checkpoint.RequestValue(req, "Hello", nil))
	t.Log(checkpoint.RequestValue(req, "encoded", nil))

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))
}
