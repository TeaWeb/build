package checkpoints

import (
	"github.com/TeaWeb/build/internal/teawaf/requests"
)

// ${requestAll}
type RequestAllCheckpoint struct {
	Checkpoint
}

func (this *RequestAllCheckpoint) RequestValue(req *requests.Request, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	valueBytes := []byte{}
	if len(req.RequestURI) > 0 {
		valueBytes = append(valueBytes, req.RequestURI...)
	} else if req.URL != nil {
		valueBytes = append(valueBytes, req.URL.RequestURI()...)
	}

	if req.Body != nil {
		valueBytes = append(valueBytes, ' ')

		if len(req.BodyData) == 0 {
			data, err := req.ReadBody(int64(32 * 1024 * 1024)) // read 32m bytes
			if err != nil {
				return "", err, nil
			}

			req.BodyData = data
			req.RestoreBody(data)
		}
		valueBytes = append(valueBytes, req.BodyData...)
	}

	value = valueBytes

	return
}

func (this *RequestAllCheckpoint) ResponseValue(req *requests.Request, resp *requests.Response, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	value = ""
	if this.IsRequest() {
		return this.RequestValue(req, param, options)
	}
	return
}
