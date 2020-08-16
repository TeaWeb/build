package apps

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
)

type ExecItemSourceAction actions.Action

// 立即执行监控项数据源
func (this *ExecItemSourceAction) Run(params struct {
	AgentId string
	AppId   string
	ItemId  string
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到要操作的Agent")
	}

	app := agent.FindApp(params.AppId)
	if app == nil {
		this.Fail("找不到App")
	}

	item := app.FindItem(params.ItemId)
	if item == nil {
		this.Fail("找不到Item")
	}

	// 通知更新
	agentutils.PostAgentEvent(agent.Id, agentutils.NewAgentEvent("RUN_ITEM", maps.Map{
		"appId":  app.Id,
		"itemId": params.ItemId,
	}))

	// 同步
	if app.IsSharedWithGroup {
		err := agentutils.SyncAppEvent(agent.Id, agent.GroupIds, app, &agentutils.AgentEvent{
			Name: "RUN_ITEM",
			Data: maps.Map{
				"appId":  app.Id,
				"itemId": params.ItemId,
			},
		})
		if err != nil {
			logs.Error(err)
		}
	}

	this.Success()
}
