package v1

import (
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/apiutils"
	"github.com/iwind/TeaGo/actions"
	"os"
	"time"
)

type StopAction actions.Action

// 停止服务
func (this *StopAction) RunGet(params struct{}) {
	go func() {
		time.Sleep(1 * time.Second)
		os.Exit(0)
	}()
	apiutils.SuccessOK(this)
}
