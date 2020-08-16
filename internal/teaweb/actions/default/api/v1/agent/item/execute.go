package item

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/agentutils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/apiutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type ExecuteAction actions.Action

// 执行一次监控项
func (this *ExecuteAction) RunGet(params struct {
	AgentId string
	AppId   string
	ItemId  string
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		apiutils.Fail(this, "agent not found")
		return
	}

	app := agent.FindApp(params.AppId)
	if app == nil {
		apiutils.Fail(this, "app not found")
		return
	}

	item := app.FindItem(params.ItemId)
	if item == nil {
		apiutils.Fail(this, "item not found")
		return
	}

	agentutils.PostAgentEvent(agent.Id, agentutils.NewAgentEvent("RUN_ITEM", maps.Map{
		"appId":  app.Id,
		"itemId": params.ItemId,
	}))

	// 同步
	if app.IsSharedWithGroup {
		agentutils.SyncAppEvent(agent.Id, agent.GroupIds, app, &agentutils.AgentEvent{
			Name: "RUN_ITEM",
			Data: maps.Map{
				"appId":  app.Id,
				"itemId": params.ItemId,
			},
		})
	}

	apiutils.SuccessOK(this)
}
