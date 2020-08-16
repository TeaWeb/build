package agentutils

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
)

func ActionDeleteAgent(agentId string, onFail func(message string)) (goNext bool) {
	agent := agents.NewAgentConfigFromId(agentId)
	if agent == nil {
		onFail("要删除的主机不存在")
		return
	}

	// 删除通知
	err := teadb.NoticeDAO().DeleteNoticesForAgent(agent.Id)
	if err != nil {
		onFail("通知删除失败：" + err.Error())
		return
	}

	// 删除数值记录
	_ = teadb.AgentValueDAO().DropAgentTable(agent.Id)
	_ = teadb.AgentLogDAO().DropAgentTable(agent.Id)

	// 从列表删除
	agentList, err := agents.SharedAgentList()
	if err != nil {
		onFail("删除失败：" + err.Error())
		return
	}
	agentList.RemoveAgent(agent.Filename())
	err = agentList.Save()
	if err != nil {
		onFail("删除失败：" + err.Error())
		return
	}

	// 删除配置文件
	err = agent.Delete()
	if err != nil {
		onFail("删除失败：" + err.Error())
		return
	}

	// 减少分组数据
	err = agents.SharedGroupList().BuildIndexes()
	if err != nil {
		logs.Error(err)
	}

	// 通知更新
	PostAgentEvent(agent.Id, NewAgentEvent("REMOVE_AGENT", maps.Map{}))
	return true
}
