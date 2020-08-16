package checkpoints

import (
	"bytes"
	"github.com/TeaWeb/build/internal/teawaf/requests"
	"github.com/iwind/TeaGo/types"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestRequestBodyCheckpoint_RequestValue(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "http://teaos.cn", bytes.NewBuffer([]byte("123456")))
	if err != nil {
		t.Fatal(err)
	}

	checkpoint := new(RequestBodyCheckpoint)
	t.Log(checkpoint.RequestValue(requests.NewRequest(req), "", nil))

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(body))
}

func TestRequestBodyCheckpoint_RequestValue_Max(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "http://teaos.cn", bytes.NewBuffer([]byte(strings.Repeat("123456", 10240000))))
	if err != nil {
		t.Fatal(err)
	}

	checkpoint := new(RequestBodyCheckpoint)
	value, err, _ := checkpoint.RequestValue(requests.NewRequest(req), "", nil)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("value bytes:", len(types.String(value)))

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("raw bytes:", len(body))
}
