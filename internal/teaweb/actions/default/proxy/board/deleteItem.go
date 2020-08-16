package board

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
)

type DeleteItemAction actions.Action

// 删除指标
func (this *DeleteItemAction) RunPost(params struct {
	ServerId string
	Code     string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到代理服务")
	}

	if len(params.Code) == 0 {
		this.Fail("请选择指标")
	}

	server.RemoveStatItem(params.Code)
	err := server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 重新启动统计
	proxyutils.ReloadServerStats(server.Id)

	this.Success()
}
