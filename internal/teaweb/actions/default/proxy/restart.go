package proxy

import (
	"github.com/TeaWeb/build/internal/teaproxy"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
)

type RestartAction actions.Action

// 试图重启所有服务
func (this *RestartAction) Run(params struct{}) {
	err := teaproxy.SharedManager.Restart()
	if err != nil {
		this.Fail("重启失败：" + err.Error())
	}

	proxyutils.FinishChange()

	this.Refresh().Success()
}
