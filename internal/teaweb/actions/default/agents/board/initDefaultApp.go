package apps

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/iwind/TeaGo/actions"
)

type InitDefaultAppAction actions.Action

// 初始化内置的App
func (this *InitDefaultAppAction) Run(params struct {
	AgentId string
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到Agent")
	}

	systemApp := agent.FindApp("system")
	if systemApp != nil {
		agent.RemoveApp(systemApp.Id)

		board := agents.NewAgentBoard(agent.Id)
		if board != nil {
			board.RemoveApp(systemApp.Id)
			err := board.Save()
			if err != nil {
				this.Fail("操作失败：" + err.Error())
			}
		}

		agent.AddDefaultApps()
	} else {
		agent.AddDefaultApps()
	}

	err := agent.Save()
	if err != nil {
		this.Fail("操作失败：" + err.Error())
	}

	this.Success()
}
