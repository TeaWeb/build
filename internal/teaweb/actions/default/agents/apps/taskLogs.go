package apps

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
)

type TaskLogsAction actions.Action

// 任务日志
func (this *TaskLogsAction) Run(params struct {
	AgentId string
	AppId   string
	TaskId  string
	Tabbar  string
}) {
	this.Data["tabbar"] = params.Tabbar

	agentutils.InitAppData(this, params.AgentId, params.AppId, params.Tabbar)

	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到Agent")
	}

	app := agent.FindApp(params.AppId)
	if app == nil {
		this.Fail("找不到App")
	}

	task := app.FindTask(params.TaskId)
	if task == nil {
		this.Fail("找不到要修改的任务")
	}

	this.Data["task"] = maps.Map{
		"id":        task.Id,
		"name":      task.Name,
		"on":        task.On,
		"script":    task.Script,
		"cwd":       task.Cwd,
		"isBooting": task.IsBooting,
		"isManual":  task.IsManual,
		"env":       task.Env,
		"schedules": lists.Map(task.Schedule, func(k int, v interface{}) interface{} {
			s := v.(*agents.ScheduleConfig)
			return maps.Map{
				"summary": s.Summary(),
			}
		}),
	}

	this.Show()
}

// 日志数据
func (this *TaskLogsAction) RunPost(params struct {
	AgentId string
	TaskId  string
	LastId  string
}) {
	taskLogs, err := teadb.AgentLogDAO().FindLatestTaskLogs(params.AgentId, params.TaskId, params.LastId, 100)
	if err != nil {
		logs.Error(err)
	}
	this.Data["logs"] = taskLogs
	this.Success()
}
