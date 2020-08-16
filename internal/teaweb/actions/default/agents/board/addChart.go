package apps

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
)

type AddChartAction actions.Action

// 添加图表
func (this *AddChartAction) Run(params struct {
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
	board.AddChart(params.AppId, params.ItemId, params.ChartId)
	err := board.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 同步
	if app.IsSharedWithGroup {
		err := agentutils.SyncAddChart(agent.Id, agent.GroupIds, app, item.Id, params.ChartId)
		if err != nil {
			logs.Error(err)
		}
	}

	this.Success()
}
