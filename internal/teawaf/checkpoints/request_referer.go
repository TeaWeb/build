package checkpoints

import (
	"github.com/TeaWeb/build/internal/teawaf/requests"
)

type RequestRefererCheckpoint struct {
	Checkpoint
}

func (this *RequestRefererCheckpoint) RequestValue(req *requests.Request, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	value = req.Referer()
	return
}

func (this *RequestRefererCheckpoint) ResponseValue(req *requests.Request, resp *requests.Response, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	if this.IsRequest() {
		return this.RequestValue(req, param, options)
	}
	return
}
