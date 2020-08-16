package apps

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
)

type ClearItemValuesAction actions.Action

// 清除数值记录
func (this *ClearItemValuesAction) Run(params struct {
	AgentId string
	AppId   string
	ItemId  string
	Level   int
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到Agent")
	}

	app := agent.FindApp(params.AppId)
	if app == nil {
		this.Fail("找不到App")
	}

	item := app.FindItem(params.ItemId)
	if item == nil {
		this.Fail("找不到Item")
	}

	err := teadb.AgentValueDAO().ClearItemValues(params.AgentId, params.AppId, params.ItemId, notices.NoticeLevel(params.Level))
	if err != nil {
		this.Fail("清除失败：" + err.Error())
	}

	// 清除同组
	if app.IsSharedWithGroup {
		for _, agent1 := range agentutils.FindSharedAgents(agent.Id, agent.GroupIds, app) {
			err := teadb.AgentValueDAO().ClearItemValues(agent1.Id, params.AppId, params.ItemId, notices.NoticeLevel(params.Level))
			if err != nil {
				this.Fail("清除失败：" + err.Error())
			}
		}
	}

	this.Success()
}
