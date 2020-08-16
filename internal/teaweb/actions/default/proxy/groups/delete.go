package groups

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
)

type DeleteAction actions.Action

// 删除分组
func (this *DeleteAction) Run(params struct {
	ServerId string
	GroupId  string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	if len(params.GroupId) == 0 {
		this.Fail("要删除的分组ID不能为空")
	}

	if params.GroupId == "default" {
		this.Fail("默认分组无法删除")
	}

	server.RemoveRequestGroup(params.GroupId)

	// 删除已经使用的分组ID
	for _, backend := range server.Backends {
		backend.RemoveRequestGroupId(params.GroupId)
	}

	for _, location := range server.Locations {
		for _, backend := range location.Backends {
			backend.RemoveRequestGroupId(params.GroupId)
		}
	}

	err := server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 通知改变
	proxyutils.NotifyChange()

	this.Success()
}
