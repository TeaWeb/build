package task

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/agentutils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/apiutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type RunAction actions.Action

// 运行任务
func (this *RunAction) RunGet(params struct {
	AgentId string
	AppId   string
	TaskId  string
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		apiutils.Fail(this, "agent not found")
		return
	}

	app := agent.FindApp(params.AppId)
	if app == nil {
		apiutils.Fail(this, "app not found")
		return
	}

	task := app.FindTask(params.TaskId)
	if task == nil {
		apiutils.Fail(this, "task not found")
		return
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

	apiutils.SuccessOK(this)
}
