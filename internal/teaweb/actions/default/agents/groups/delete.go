package groups

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/iwind/TeaGo/actions"
)

type DeleteAction actions.Action

// 删除分组
func (this *DeleteAction) Run(params struct {
	GroupId string
}) {
	if len(params.GroupId) == 0 {
		this.Fail("请输入要删除的分组ID")
	}

	if params.GroupId == "default" {
		this.Fail("无法删除默认分组")
	}

	// 删除agent中的groupId
	agentList, err := agents.SharedAgentList()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	for _, agent := range agentList.FindAllAgents() {
		agent.RemoveGroup(params.GroupId)
		err = agent.Save()
		if err != nil {
			this.Fail("保存失败：" + err.Error())
		}
	}

	config := agents.SharedGroupList()
	config.RemoveGroup(params.GroupId)
	err = config.Save()
	if err != nil {
		this.Fail("保存失败： " + err.Error())
	}

	this.Success()
}
