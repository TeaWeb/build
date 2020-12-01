package apps

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/rands"
)

type AddDefaultChartsAction actions.Action

// 添加数据源默认图表
func (this *AddDefaultChartsAction) Run(params struct {
	AgentId string
	AppId   string
	ItemId  string
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到要修改的Agent")
	}

	app := agent.FindApp(params.AppId)
	if app == nil {
		this.Fail("找不到要修改的App")
	}

	item := app.FindItem(params.ItemId)
	if item == nil {
		this.Fail("找不到要操作的Item")
	}

	source := item.Source()
	if source == nil {
		this.Fail("找不到数据源")
	}

	for _, c := range source.Charts() {
		c.Id = rands.HexString(16)
		item.AddChart(c)
	}

	err := agent.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 同步
	if app.IsSharedWithGroup {
		agentutils.SyncApp(agent.Id, agent.GroupIds, app, nil, nil)
	}

	this.Success()
}
