package checkpoints

import (
	"github.com/TeaWeb/build/internal/teawaf/requests"
)

type RequestHostCheckpoint struct {
	Checkpoint
}

func (this *RequestHostCheckpoint) RequestValue(req *requests.Request, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	value = req.Host
	return
}

func (this *RequestHostCheckpoint) ResponseValue(req *requests.Request, resp *requests.Response, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	if this.IsRequest() {
		return this.RequestValue(req, param, options)
	}
	return
}
