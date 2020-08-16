package apps

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
)

type DeleteItemChartAction actions.Action

// 删除Chart
func (this *DeleteItemChartAction) Run(params struct {
	AgentId string
	AppId   string
	ItemId  string
	ChartId string
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

	// 从看板中删除
	board := agents.NewAgentBoard(params.AgentId)
	if board != nil {
		board.RemoveChart(params.ChartId)
		err := board.Save()
		if err != nil {
			this.Fail("删除失败：" + err.Error())
		}
	}

	item.RemoveChart(params.ChartId)
	err := agent.Save()
	if err != nil {
		this.Fail("删除失败：" + err.Error())
	}

	// 同步
	if app.IsSharedWithGroup {
		agentutils.SyncApp(agent.Id, agent.GroupIds, app, nil, func(agent *agents.AgentConfig) error {
			// 从看板中删除
			board := agents.NewAgentBoard(agent.Id)
			if board != nil {
				board.RemoveChart(params.ChartId)
				err := board.Save()
				if err != nil {
					this.Fail("删除失败：" + err.Error())
				}
			}
			return nil
		})
	}

	this.Success()
}
