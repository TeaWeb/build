package checkpoints

import (
	"github.com/TeaWeb/build/internal/teawaf/requests"
)

// ${requestBody}
type RequestBodyCheckpoint struct {
	Checkpoint
}

func (this *RequestBodyCheckpoint) RequestValue(req *requests.Request, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	if req.Body == nil {
		value = ""
		return
	}

	if len(req.BodyData) == 0 {
		data, err := req.ReadBody(int64(32 * 1024 * 1024)) // read 32m bytes
		if err != nil {
			return "", err, nil
		}

		req.BodyData = data
		req.RestoreBody(data)
	}

	return req.BodyData, nil, nil
}

func (this *RequestBodyCheckpoint) ResponseValue(req *requests.Request, resp *requests.Response, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	if this.IsRequest() {
		return this.RequestValue(req, param, options)
	}
	return
}
