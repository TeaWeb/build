package apps

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type OnAction actions.Action

// 启用App
func (this *OnAction) Run(params struct {
	AgentId string
	AppId   string
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到Agent")
	}

	app := agent.FindApp(params.AppId)
	if app == nil {
		this.Fail("找不到App")
	}

	app.On = true
	err := agent.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 通知更新
	agentutils.PostAgentEvent(agent.Id, agentutils.NewAgentEvent("UPDATE_APP", maps.Map{
		"appId": app.Id,
	}))

	// 同步
	if app.IsSharedWithGroup {
		agentutils.SyncApp(agent.Id, agent.GroupIds, app, agentutils.NewAgentEvent("UPDATE_APP", maps.Map{
			"appId": app.Id,
		}), nil)
	}

	this.Success()
}
