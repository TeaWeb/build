package agents

import (
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
)

type DeleteAction actions.Action

// 删除
func (this *DeleteAction) Run(params struct {
	AgentId string
}) {
	this.Data["agentId"] = params.AgentId

	if params.AgentId == "local" {
		this.RedirectURL("/agents/board")
		return
	}

	this.Show()
}

// 提交
func (this *DeleteAction) RunPost(params struct {
	AgentId string
}) {
	if !agentutils.ActionDeleteAgent(params.AgentId, func(message string) {
		this.Fail(message)
	}) {
		return
	}

	this.Success()
}
