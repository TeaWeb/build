package tunnel

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaproxy"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction actions.Action

// tunnel设置
func (this *IndexAction) RunGet(params struct {
	ServerId string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	proxyutils.AddServerMenu(this, true)

	this.Data["selectedTab"] = "tunnel"
	this.Data["server"] = proxyutils.WrapServerData(server)
	this.Data["tunnel"] = server.Tunnel

	// 状态
	runningServer := teaproxy.SharedManager.FindServer(server.Id)
	this.Data["isActive"] = false
	this.Data["errors"] = []string{}
	this.Data["countConnections"] = 0
	if runningServer != nil && runningServer.Tunnel != nil {
		this.Data["isActive"] = runningServer.Tunnel.IsActive()
		this.Data["errors"] = runningServer.Tunnel.Errors()

		tunnel := teaproxy.SharedTunnelManager.FindTunnel(server.Id, runningServer.Tunnel.Id)
		if tunnel != nil {
			this.Data["countConnections"] = tunnel.CountConnections()
		} else {
			this.Data["countConnections"] = 0
		}
	}

	this.Show()
}
