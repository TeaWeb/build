package ssl

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaproxy"
	"github.com/iwind/TeaGo/actions"
)

type ShutdownHttpsAction actions.Action

// 启动
func (this *ShutdownHttpsAction) Run(params struct {
	ServerId string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	if server.SSL == nil {
		this.Fail("还没有配置HTTPS")
	}

	if server.SSL.On {
		server.SSL.On = false
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
