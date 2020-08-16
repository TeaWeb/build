package tunnel

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
)

type UpdateAction actions.Action

// 修改
func (this *UpdateAction) RunGet(params struct {
	ServerId string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	proxyutils.AddServerMenu(this, true)

	this.Data["selectedTab"] = "tunnel"
	this.Data["server"] = proxyutils.WrapServerData(server)
	this.Data["hasTunnel"] = server.Tunnel != nil

	if server.Tunnel != nil {
		this.Data["tunnel"] = server.Tunnel
	} else {
		this.Data["tunnel"] = &teaconfigs.TunnelConfig{
			On: true,
		}
	}

	this.Show()
}

// 提交保存
func (this *UpdateAction) RunPost(params struct {
	ServerId string
	Endpoint string
	Secret   string
	On       bool
	Must     *actions.Must
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	params.Must.
		Field("endpoint", params.Endpoint).
		Require("请输入服务器终端地址")

	server.Tunnel = &teaconfigs.TunnelConfig{
		On:       params.On,
		Endpoint: teautils.FormatAddress(params.Endpoint),
		Secret:   params.Secret,
	}

	err := server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 通知更新
	proxyutils.NotifyChange()

	this.Success()
}
