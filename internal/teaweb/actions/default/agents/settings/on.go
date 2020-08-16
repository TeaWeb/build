package settings

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/iwind/TeaGo/actions"
)

type OnAction actions.Action

// 启用Agent
func (this *OnAction) Run(params struct {
	AgentId string
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到Agent")
	}

	agent.On = true
	err := agent.Save()
	if err != nil {
		this.Fail("保存失败:" + err.Error())
	}

	this.Success()
}
