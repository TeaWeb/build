package proxy

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/actions"
)

type FrontendAction actions.Action

// 前端设置
func (this *FrontendAction) Run(params struct {
	ServerId string
}) {
	this.Data["selectedTab"] = "frontend"

	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	this.Data["server"] = server

	this.Show()
}
