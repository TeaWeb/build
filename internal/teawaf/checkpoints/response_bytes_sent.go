package checkpoints

import (
	"github.com/TeaWeb/build/internal/teawaf/requests"
)

// ${bytesSent}
type ResponseBytesSentCheckpoint struct {
	Checkpoint
}

func (this *ResponseBytesSentCheckpoint) IsRequest() bool {
	return false
}

func (this *ResponseBytesSentCheckpoint) RequestValue(req *requests.Request, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	value = 0
	return
}

func (this *ResponseBytesSentCheckpoint) ResponseValue(req *requests.Request, resp *requests.Response, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	value = 0
	if resp != nil {
		value = resp.ContentLength
	}
	return
}
