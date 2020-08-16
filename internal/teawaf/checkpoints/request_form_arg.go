package checkpoints

import (
	"github.com/TeaWeb/build/internal/teawaf/requests"
	"net/url"
)

// ${requestForm.arg}
type RequestFormArgCheckpoint struct {
	Checkpoint
}

func (this *RequestFormArgCheckpoint) RequestValue(req *requests.Request, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	if req.Body == nil {
		value = ""
		return
	}

	if len(req.BodyData) == 0 {
		data, err := req.ReadBody(32 * 1024 * 1024) // read 32m bytes
		if err != nil {
			return "", err, nil
		}

		req.BodyData = data
		req.RestoreBody(data)
	}

	// TODO improve performance
	values, _ := url.ParseQuery(string(req.BodyData))
	return values.Get(param), nil, nil
}

func (this *RequestFormArgCheckpoint) ResponseValue(req *requests.Request, resp *requests.Response, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	if this.IsRequest() {
		return this.RequestValue(req, param, options)
	}
	return
}
