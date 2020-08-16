package agents

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
)

type DeleteAgentsAction actions.Action

// 删除一组Agents
func (this *DeleteAgentsAction) Run(params struct {
	AgentIds []string
}) {
	for _, agentId := range params.AgentIds {
		// 跳过本地主机
		if agentId == "local" {
			continue
		}
		agent := agents.NewAgentConfigFromId(agentId)
		if agent == nil {
			continue
		}

		// 删除通知
		err := teadb.NoticeDAO().DeleteNoticesForAgent(agent.Id)
		if err != nil {
			this.Fail("通知删除失败：" + err.Error())
		}

		// 删除数值记录
		_ = teadb.AgentValueDAO().DropAgentTable(agent.Id)
		_ = teadb.AgentLogDAO().DropAgentTable(agent.Id)

		// 从列表删除
		agentList, err := agents.SharedAgentList()
		if err != nil {
			logs.Error(err)
			continue
		}
		agentList.RemoveAgent(agent.Filename())
		err = agentList.Save()
		if err != nil {
			logs.Error(err)
			continue
		}

		err = agent.Delete()
		if err != nil {
			logs.Error(err)
			continue
		}

		// 删除通知
		err = teadb.NoticeDAO().DeleteNoticesForAgent(agent.Id)
		if err != nil {
			logs.Error(err)
		}

		// 通知更新
		agentutils.PostAgentEvent(agent.Id, agentutils.NewAgentEvent("REMOVE_AGENT", maps.Map{}))
	}

	// 重建索引
	err := agents.SharedGroupList().BuildIndexes()
	if err != nil {
		logs.Error(err)
	}

	this.Success()
}
