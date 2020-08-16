package proxy

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaproxy"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
)

type DetailAction actions.Action

// 代理详情
func (this *DetailAction) Run(params struct {
	ServerId string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	if server.Index == nil {
		server.Index = []string{}
	}

	this.Data["selectedTab"] = "basic"
	this.Data["server"] = server
	this.Data["isTCP"] = server.IsTCP()
	this.Data["isForwardHTTP"] = server.IsForwardHTTP()

	this.Data["errs"] = teaproxy.SharedManager.FindServerErrors(params.ServerId)
	this.Data["accessLogs"] = proxyutils.FormatAccessLog(server.AccessLog)

	this.Show()
}
