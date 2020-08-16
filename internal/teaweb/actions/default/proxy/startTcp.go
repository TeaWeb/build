package proxy

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaproxy"
	"github.com/iwind/TeaGo/actions"
)

type StartTcpAction actions.Action

// 启动Tcp服务
func (this *StartTcpAction) Run(params struct {
	ServerId string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	isChanged := false
	if !server.On {
		server.On = true
		isChanged = true
	}

	if server.IsTCP() {
		server.TCP.TCPOn = true
		isChanged = true
	}

	if !isChanged {
		this.Success()
	}

	err := server.Validate()
	if err != nil {
		this.Fail("Server校验失败：" + err.Error())
	}

	err = server.Save()
	if err != nil {
		this.Fail("启动失败：" + err.Error())
	}

	teaproxy.SharedManager.ApplyServer(server)
	teaproxy.SharedManager.Reload()

	this.Success()
}
