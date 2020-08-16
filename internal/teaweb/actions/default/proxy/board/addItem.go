package board

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teastats"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
)

type AddItemAction actions.Action

// 添加数据指标
func (this *AddItemAction) RunPost(params struct {
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

	filter := teastats.FindSharedFilter(params.Code)
	if filter == nil {
		this.Fail("请选择指标")
	}

	server.AddStatItem(params.Code)
	err := server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 重新启动统计
	proxyutils.ReloadServerStats(server.Id)

	this.Success()
}
