package backend

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction actions.Action

// 后端列表
func (this *IndexAction) Run(params struct {
	ServerId string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	this.Data["selectedTab"] = "backend"
	this.Data["server"] = server
	this.Data["location"] = nil
	this.Data["isTCP"] = server.IsTCP()
	this.Data["isHTTP"] = server.IsHTTP()

	this.Data["queryParams"] = maps.Map{
		"serverId": params.ServerId,
	}

	this.Show()
}
