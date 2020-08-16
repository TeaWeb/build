package backend

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaproxy"
	"github.com/iwind/TeaGo/actions"
)

type ClearFailsAction actions.Action

// 清除失败次数
func (this *ClearFailsAction) Run(params struct {
	ServerId   string
	LocationId string
	Websocket  bool
	BackendId  string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	runningServer := teaproxy.SharedManager.FindServer(server.Id)
	if runningServer != nil {
		backendList, _ := runningServer.FindBackendList(params.LocationId, params.Websocket)
		if backendList != nil {
			backend := backendList.FindBackend(params.BackendId)
			if backend != nil {
				backend.CurrentFails = 0
			}
		}
	}

	this.Success()
}
