package v1

import (
	"github.com/TeaWeb/build/internal/teaproxy"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/apiutils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
)

type ReloadAction actions.Action

// 重载代理服务
func (this *ReloadAction) RunGet(params struct{}) {
	err := teaproxy.SharedManager.Restart()
	if err != nil {
		apiutils.Fail(this, err.Error())
		return
	}

	proxyutils.FinishChange()

	apiutils.SuccessOK(this)
}
