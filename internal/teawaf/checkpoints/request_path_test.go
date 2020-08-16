package checkpoints

import (
	"github.com/TeaWeb/build/internal/teawaf/requests"
	"net/http"
	"testing"
)

func TestRequestPathCheckpoint_RequestValue(t *testing.T) {
	rawReq, err := http.NewRequest(http.MethodGet, "http://teaos.cn/index?name=lu", nil)
	if err != nil {
		t.Fatal(err)
	}

	req := requests.NewRequest(rawReq)
	checkpoint := new(RequestPathCheckpoint)
	t.Log(checkpoint.RequestValue(req, "", nil))
}
