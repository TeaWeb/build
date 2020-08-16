package proxy

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaproxy"
	"github.com/iwind/TeaGo/actions"
)

type ShutdownHttpAction actions.Action

// 关闭HTTP服务
func (this *ShutdownHttpAction) Run(params struct {
	ServerId string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	if server.Http {
		server.Http = false
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
