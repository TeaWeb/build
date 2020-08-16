package checkpoints

import (
	"github.com/TeaWeb/build/internal/teawaf/requests"
	"net/http"
	"testing"
)

func TestCCCheckpoint_RequestValue(t *testing.T) {
	raw, err := http.NewRequest(http.MethodGet, "http://teaos.cn/", nil)
	if err != nil {
		t.Fatal(err)
	}

	req := requests.NewRequest(raw)
	req.RemoteAddr = "127.0.0.1"

	checkpoint := new(CCCheckpoint)
	checkpoint.Init()
	checkpoint.Start()

	options := map[string]string{
		"period": "5",
	}
	t.Log(checkpoint.RequestValue(req, "requests", options))
	t.Log(checkpoint.RequestValue(req, "requests", options))

	req.RemoteAddr = "127.0.0.2"
	t.Log(checkpoint.RequestValue(req, "requests", options))

	req.RemoteAddr = "127.0.0.1"
	t.Log(checkpoint.RequestValue(req, "requests", options))

	req.RemoteAddr = "127.0.0.2"
	t.Log(checkpoint.RequestValue(req, "requests", options))

	req.RemoteAddr = "127.0.0.2"
	t.Log(checkpoint.RequestValue(req, "requests", options))

	req.RemoteAddr = "127.0.0.2"
	t.Log(checkpoint.RequestValue(req, "requests", options))
}
