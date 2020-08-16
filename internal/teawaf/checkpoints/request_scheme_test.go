package checkpoints

import (
	"github.com/TeaWeb/build/internal/teawaf/requests"
	"net/http"
	"testing"
)

func TestRequestSchemeCheckpoint_RequestValue(t *testing.T) {
	rawReq, err := http.NewRequest(http.MethodGet, "https://teaos.cn/?name=lu", nil)
	if err != nil {
		t.Fatal(err)
	}

	req := requests.NewRequest(rawReq)
	checkpoint := new(RequestSchemeCheckpoint)
	t.Log(checkpoint.RequestValue(req, "", nil))
}
