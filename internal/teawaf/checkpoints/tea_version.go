package checkpoints

import (
	"github.com/TeaWeb/build/internal/teaconst"
	"github.com/TeaWeb/build/internal/teawaf/requests"
)

type TeaVersionCheckpoint struct {
	Checkpoint
}

func (this *TeaVersionCheckpoint) RequestValue(requests *requests.Request, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	value = teaconst.TeaVersion
	return
}

func (this *TeaVersionCheckpoint) ResponseValue(requests *requests.Request, resp *requests.Response, param string, options map[string]string) (value interface{}, sysErr error, userErr error) {
	value = teaconst.TeaVersion
	return
}
