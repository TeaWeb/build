package checkpoints

import (
	"github.com/TeaWeb/build/internal/teawaf/requests"
)

// just a sample checkpoint, copy and change it for your new checkpoint
type SampleResponseCheckpoint struct {
	Checkpoint
}

func (this *SampleResponseCheckpoint) IsRequest() bool {
	return false
}

func (this *SampleResponseCheckpoint) RequestValue(req *requests.Request, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	return
}

func (this *SampleResponseCheckpoint) ResponseValue(req *requests.Request, resp *requests.Response, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	return
}
