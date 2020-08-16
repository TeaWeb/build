package apps

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type AddAction actions.Action

// 添加App
func (this *AddAction) Run(params struct {
	AgentId string
}) {
	this.Data["agentId"] = params.AgentId
	this.Show()
}

// 提交保存
func (this *AddAction) RunPost(params struct {
	AgentId           string
	Name              string
	On                bool
	IsSharedWithGroup bool
	Must              *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入App名称")

	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到Agent")
	}

	app := agents.NewAppConfig()
	app.On = params.On
	app.Name = params.Name
	app.IsSharedWithGroup = params.IsSharedWithGroup
	agent.AddApp(app)
	err := agent.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 通知更新
	agentutils.PostAgentEvent(agent.Id, agentutils.NewAgentEvent("ADD_APP", maps.Map{
		"appId": app.Id,
	}))

	// 同步
	if app.IsSharedWithGroup {
		agentutils.SyncApp(agent.Id, agent.GroupIds, app, agentutils.NewAgentEvent("ADD_APP", maps.Map{
			"appId": app.Id,
		}), nil)
	}

	this.Success()
}
