package apps

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type RunTaskAction actions.Action

// 运行一次任务
func (this *RunTaskAction) Run(params struct {
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

	agentutils.PostAgentEvent(params.AgentId, &agentutils.AgentEvent{
		Name: "RUN_TASK",
		Data: maps.Map{
			"taskId": params.TaskId,
		},
	})

	// 同步
	if app.IsSharedWithGroup {
		agentutils.SyncAppEvent(agent.Id, agent.GroupIds, app, &agentutils.AgentEvent{
			Name: "RUN_TASK",
			Data: maps.Map{
				"taskId": params.TaskId,
			},
		})
	}

	this.Success()
}
