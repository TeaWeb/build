package apps

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
)

type RemoveChartAction actions.Action

// 移除图表
func (this *RemoveChartAction) Run(params struct {
	AgentId string
	AppId   string
	ItemId  string
	ChartId string
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到Agent")
	}

	app := agent.FindApp(params.AppId)
	if app == nil {
		this.Fail("找不到App")
	}

	item := app.FindItem(params.ItemId)
	if item == nil {
		this.Fail("找不到Item")
	}

	chart := item.FindChart(params.ChartId)
	if chart == nil {
		this.Fail("找不到Chart")
	}

	board := agents.NewAgentBoard(params.AgentId)
	if board == nil {
		this.Fail("无法读取Board配置文件")
	}
	board.RemoveChart(params.ChartId)
	err := board.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 同步
	if app.IsSharedWithGroup {
		agentutils.SyncRemoveChart(agent.Id, agent.GroupIds, app, params.ChartId)
	}

	this.Success()
}
