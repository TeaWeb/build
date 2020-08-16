package locations

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
)

type MoveDownAction actions.Action

func (this *MoveDownAction) Run(params struct {
	ServerId string
	Index    int
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	if params.Index >= 0 && params.Index < len(server.Locations)-1 {
		next := server.Locations[params.Index+1]
		current := server.Locations[params.Index]
		server.Locations[params.Index+1] = current
		server.Locations[params.Index] = next
	}

	err := server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	proxyutils.NotifyChange()

	this.Refresh().Success()
}
