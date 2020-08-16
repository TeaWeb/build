package apps

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type DeleteItemAction actions.Action

// 删除监控项
func (this *DeleteItemAction) Run(params struct {
	AgentId string
	AppId   string
	ItemId  string
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

	// 删除看板中相关图表
	if len(item.Charts) > 0 {
		board := agents.NewAgentBoard(params.AgentId)
		if board != nil {
			for _, c := range item.Charts {
				board.RemoveChart(c.Id)
			}
			err := board.Save()
			if err != nil {
				this.Fail("删除失败：" + err.Error())
			}
		}
	}

	app.RemoveItem(params.ItemId)
	err := agent.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 通知更新
	agentutils.PostAgentEvent(agent.Id, agentutils.NewAgentEvent("DELETE_ITEM", maps.Map{
		"appId":  app.Id,
		"itemId": params.ItemId,
	}))

	// 同步App
	if app.IsSharedWithGroup {
		err := agentutils.SyncApp(agent.Id, agent.GroupIds, app, agentutils.NewAgentEvent("DELETE_ITEM", maps.Map{
			"appId":  app.Id,
			"itemId": params.ItemId,
		}), func(agent *agents.AgentConfig) error {
			// 删除看板中相关图表
			if len(item.Charts) > 0 {
				board := agents.NewAgentBoard(agent.Id)
				if board != nil {
					for _, c := range item.Charts {
						board.RemoveChart(c.Id)
					}
					err := board.Save()
					if err != nil {
						this.Fail("删除失败：" + err.Error())
					}
				}
			}
			return nil
		})
		if err != nil {
			this.Fail("保存失败：" + err.Error())
		}
	}

	this.Success()
}
