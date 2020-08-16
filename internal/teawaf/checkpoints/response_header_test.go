package checkpoints

import (
	"github.com/TeaWeb/build/internal/teawaf/requests"
	"net/http"
	"testing"
)

func TestResponseHeaderCheckpoint_ResponseValue(t *testing.T) {
	resp := requests.NewResponse(new(http.Response))
	resp.StatusCode = 200
	resp.Header = http.Header{}
	resp.Header.Set("Hello", "World")

	checkpoint := new(ResponseHeaderCheckpoint)
	t.Log(checkpoint.ResponseValue(nil, resp, "Hello", nil))
}
