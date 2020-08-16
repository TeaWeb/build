package apps

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/iwind/TeaGo/actions"
)

type MoveAction actions.Action

// 移动App位置
func (this *MoveAction) RunPost(params struct {
	AgentId   string
	FromIndex int
	ToIndex   int
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到Agent")
	}

	agent.MoveApp(params.FromIndex, params.ToIndex)
	err := agent.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
