package board

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaproxy"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction actions.Action

// 看板
func (this *IndexAction) Run(params struct {
	ServerId string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到要查看的代理服务")
	}

	this.Data["server"] = maps.Map{
		"id": server.Id,
	}

	this.Data["boardType"] = "stat"

	this.Data["errs"] = teaproxy.SharedManager.FindServerErrors(params.ServerId)

	this.View("/proxy/board/index.html")
	this.Show()
}
