package headers

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction actions.Action

// 自定义Http Header
func (this *IndexAction) Run(params struct {
	ServerId string // 必填
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	this.Data["selectedTab"] = "header"
	this.Data["server"] = server

	this.Data["headerQuery"] = maps.Map{
		"serverId": params.ServerId,
	}

	this.Show()
}
