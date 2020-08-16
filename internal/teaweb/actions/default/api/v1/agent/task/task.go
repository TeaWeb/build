package task

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/apiutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type TaskAction actions.Action

// 任务信息
func (this *TaskAction) RunGet(params struct {
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

	apiutils.Success(this, maps.Map{
		"config": task,
	})
}
