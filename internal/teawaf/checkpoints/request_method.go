package checkpoints

import (
	"github.com/TeaWeb/build/internal/teawaf/requests"
)

type RequestMethodCheckpoint struct {
	Checkpoint
}

func (this *RequestMethodCheckpoint) RequestValue(req *requests.Request, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	value = req.Method
	return
}

func (this *RequestMethodCheckpoint) ResponseValue(req *requests.Request, resp *requests.Response, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	if this.IsRequest() {
		return this.RequestValue(req, param, options)
	}
	return
}
