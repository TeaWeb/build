package checkpoints

import (
	"github.com/TeaWeb/build/internal/teawaf/requests"
	"net/http"
	"testing"
)

func TestRequestHostCheckpoint_RequestValue(t *testing.T) {
	rawReq, err := http.NewRequest(http.MethodGet, "https://teaos.cn/?name=lu", nil)
	if err != nil {
		t.Fatal(err)
	}

	req := requests.NewRequest(rawReq)
	req.Header.Set("Host", "cloud.teaos.cn")

	checkpoint := new(RequestHostCheckpoint)
	t.Log(checkpoint.RequestValue(req, "", nil))
}
