package apps

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type DeleteTaskAction actions.Action

// 删除任务
func (this *DeleteTaskAction) Run(params struct {
	AgentId string
	AppId   string
	TaskId  string
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到Agent")
	}

	app := agent.FindApp(params.AppId)
	if app == nil {
		this.Fail("找不到App")
	}
	app.RemoveTask(params.TaskId)

	err := agent.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 通知更新
	agentutils.PostAgentEvent(agent.Id, agentutils.NewAgentEvent("REMOVE_TASK", maps.Map{
		"appId":  app.Id,
		"taskId": params.TaskId,
	}))

	// 同步
	if app.IsSharedWithGroup {
		agentutils.SyncApp(agent.Id, agent.GroupIds, app, agentutils.NewAgentEvent("REMOVE_TASK", maps.Map{
			"appId":  app.Id,
			"taskId": params.TaskId,
		}), nil)
	}

	this.Success()
}
