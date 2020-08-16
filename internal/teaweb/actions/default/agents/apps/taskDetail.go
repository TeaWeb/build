package apps

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/time"
	"time"
)

type TaskDetailAction actions.Action

// 任务详情
func (this *TaskDetailAction) Run(params struct {
	AgentId string
	AppId   string
	TaskId  string
	From    string
	Tabbar  string
}) {
	this.Data["from"] = params.From
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

	// 下次运行时间
	nextTimeString := ""
	if len(task.Schedule) > 0 {
		err := task.Validate()
		if err == nil {
			nextTime, ok := task.Next(time.Now())
			if ok {
				nextTimeString = timeutil.Format("Y-m-d H:i:s", nextTime)
			}
		}
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
		"nextTime": nextTimeString,
	}

	this.Show()
}
