package settings

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/iwind/TeaGo/actions"
)

type OffAction actions.Action

// 关闭Agent
func (this *OffAction) Run(params struct {
	AgentId string
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到Agent")
	}

	agent.On = false
	err := agent.Save()
	if err != nil {
		this.Fail("保存失败:" + err.Error())
	}

	this.Success()
}
