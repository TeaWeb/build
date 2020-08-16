package apps

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type UpdateAction actions.Action

// 修改
func (this *UpdateAction) Run(params struct {
	From    string
	AgentId string
	AppId   string
}) {
	this.Data["from"] = params.From

	agentutils.InitAppData(this, params.AgentId, params.AppId, "detail")
	this.Show()
}

// 提交保存
func (this *UpdateAction) RunPost(params struct {
	AgentId           string
	AppId             string
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

	app := agent.FindApp(params.AppId)
	if app == nil {
		this.Fail("找不到要修改的App")
	}

	isSharedWithGroup := app.IsSharedWithGroup

	app.On = params.On
	app.Name = params.Name
	app.IsSharedWithGroup = params.IsSharedWithGroup
	err := agent.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 通知更新
	agentutils.PostAgentEvent(agent.Id, agentutils.NewAgentEvent("UPDATE_APP", maps.Map{
		"appId": app.Id,
	}))

	// 同步
	if isSharedWithGroup || app.IsSharedWithGroup {
		agentutils.SyncApp(agent.Id, agent.GroupIds, app, agentutils.NewAgentEvent("UPDATE_APP", maps.Map{
			"appId": app.Id,
		}), nil)
	}

	this.Success()
}
