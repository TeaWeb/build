package checkpoints

import (
	"github.com/TeaWeb/build/internal/teawaf/requests"
)

type RequestURICheckpoint struct {
	Checkpoint
}

func (this *RequestURICheckpoint) RequestValue(req *requests.Request, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	if len(req.RequestURI) > 0 {
		value = req.RequestURI
	} else if req.URL != nil {
		value = req.URL.RequestURI()
	}
	return
}

func (this *RequestURICheckpoint) ResponseValue(req *requests.Request, resp *requests.Response, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	if this.IsRequest() {
		return this.RequestValue(req, param, options)
	}
	return
}
