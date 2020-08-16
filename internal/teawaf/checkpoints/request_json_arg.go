package checkpoints

import (
	"encoding/json"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/TeaWeb/build/internal/teawaf/requests"
	"strings"
)

// ${requestJSON.arg}
type RequestJSONArgCheckpoint struct {
	Checkpoint
}

func (this *RequestJSONArgCheckpoint) RequestValue(req *requests.Request, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	if len(req.BodyData) == 0 {
		data, err := req.ReadBody(int64(32 * 1024 * 1024)) // read 32m bytes
		if err != nil {
			return "", err, nil
		}
		req.BodyData = data
		defer req.RestoreBody(data)
	}

	// TODO improve performance
	var m interface{} = nil
	err := json.Unmarshal(req.BodyData, &m)
	if err != nil || m == nil {
		return "", nil, err
	}

	value = teautils.Get(m, strings.Split(param, "."))
	if value != nil {
		return value, nil, err
	}
	return "", nil, nil
}

func (this *RequestJSONArgCheckpoint) ResponseValue(req *requests.Request, resp *requests.Response, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	if this.IsRequest() {
		return this.RequestValue(req, param, options)
	}
	return
}
