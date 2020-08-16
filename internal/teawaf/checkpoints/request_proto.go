package checkpoints

import (
	"github.com/TeaWeb/build/internal/teawaf/requests"
)

type RequestProtoCheckpoint struct {
	Checkpoint
}

func (this *RequestProtoCheckpoint) RequestValue(req *requests.Request, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	value = req.Proto
	return
}

func (this *RequestProtoCheckpoint) ResponseValue(req *requests.Request, resp *requests.Response, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	if this.IsRequest() {
		return this.RequestValue(req, param, options)
	}
	return
}
