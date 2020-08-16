package checkpoints

import (
	"github.com/TeaWeb/build/internal/teawaf/requests"
)

type RequestSchemeCheckpoint struct {
	Checkpoint
}

func (this *RequestSchemeCheckpoint) RequestValue(req *requests.Request, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	value = req.URL.Scheme
	return
}

func (this *RequestSchemeCheckpoint) ResponseValue(req *requests.Request, resp *requests.Response, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	if this.IsRequest() {
		return this.RequestValue(req, param, options)
	}
	return
}
