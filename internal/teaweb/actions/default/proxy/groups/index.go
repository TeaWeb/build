package groups

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/actions"
)

type IndexAction actions.Action

// 分组管理
func (this *IndexAction) Run(params struct {
	ServerId   string
	LocationId string
	Websocket  int
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	if len(params.LocationId) > 0 {
		this.Data["selectedTab"] = "location"
	} else {
		this.Data["selectedTab"] = "backend"
	}
	this.Data["server"] = server
	this.Data["locationId"] = params.LocationId
	this.Data["websocket"] = params.Websocket

	if len(server.RequestGroups) > 0 {
		this.Data["groups"] = server.RequestGroups
	} else {
		this.Data["groups"] = []*teaconfigs.RequestGroup{}
	}

	this.Show()
}
