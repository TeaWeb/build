package checkpoints

import (
	"github.com/TeaWeb/build/internal/teawaf/requests"
	"github.com/iwind/TeaGo/types"
	"net"
)

type RequestRemotePortCheckpoint struct {
	Checkpoint
}

func (this *RequestRemotePortCheckpoint) RequestValue(req *requests.Request, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	_, port, err := net.SplitHostPort(req.RemoteAddr)
	if err == nil {
		value = types.Int(port)
	} else {
		value = 0
	}
	return
}

func (this *RequestRemotePortCheckpoint) ResponseValue(req *requests.Request, resp *requests.Response, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	if this.IsRequest() {
		return this.RequestValue(req, param, options)
	}
	return
}
