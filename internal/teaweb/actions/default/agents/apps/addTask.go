package apps

import (
	"encoding/json"
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaconst"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"reflect"
)

type AddTaskAction actions.Action

// 添加任务
func (this *AddTaskAction) Run(params struct {
	From         string
	AgentId      string
	AppId        string
	IsScheduling bool
	IsBooting    bool
	IsManual     bool
}) {
	this.Data["from"] = params.From

	tabbar := ""
	if params.IsScheduling {
		tabbar = "schedule"
	} else if params.IsBooting {
		tabbar = "boot"
	} else if params.IsManual {
		tabbar = "manual"
	}
	this.Data["tabbar"] = tabbar
	agentutils.InitAppData(this, params.AgentId, params.AppId, tabbar)

	this.Show()
}

func (this *AddTaskAction) RunPost(params struct {
	AgentId       string
	AppId         string
	Name          string
	Script        string
	Cwd           string
	EnvNames      []string
	EnvValues     []string
	SchedulesJSON string
	IsBooting     bool
	IsManual      bool
	On            bool
	Must          *actions.Must
}) {
	if teaconst.DemoEnabled {
		this.Fail("演示版无法添加任务")
	}

	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到Agent")
	}

	app := agent.FindApp(params.AppId)
	if app == nil {
		this.Fail("找不到App")
	}

	params.Must.
		Field("name", params.Name).
		Require("请输入任务名").
		Field("script", params.Script).
		Require("请输入Shell脚本内容")

	task := agents.NewTaskConfig()
	task.On = params.On
	task.Name = params.Name
	task.Script = params.Script
	task.Cwd = params.Cwd
	for index, name := range params.EnvNames {
		if index < len(params.EnvValues) {
			task.AddEnv(name, params.EnvValues[index])
		}
	}
	task.IsBooting = params.IsBooting
	task.IsManual = params.IsManual

	if len(params.SchedulesJSON) > 0 {
		rangeArray := []map[string]interface{}{}
		err := json.Unmarshal([]byte(params.SchedulesJSON), &rangeArray)
		if err != nil {
			this.Fail("定时任务设置失败：" + err.Error())
		}

		for _, m := range rangeArray {
			schedule := agents.NewScheduleConfig()
			for timeType, timeConfig := range m {
				timeTypeString := types.String(timeType)
				timeConfigMap := maps.NewMap(timeConfig)
				ranges := []*agents.ScheduleRangeConfig{}
				{
					every := timeConfigMap.GetBool("every")
					if every {
						r := agents.NewScheduleRangeConfig()
						r.Every = every
						ranges = append(ranges, r)
					}
				}
				{
					points := timeConfigMap.Get("points")
					if points != nil && reflect.TypeOf(points).Kind() == reflect.Slice {
						lists.NewList(points).Range(func(k int, v interface{}) {
							i := types.Int(v)
							if i >= 0 {
								r := agents.NewScheduleRangeConfig()
								r.Value = i
								ranges = append(ranges, r)
							}
						})
					}
				}
				{
					steps := timeConfigMap.Get("steps")
					if steps != nil && reflect.TypeOf(steps).Kind() == reflect.Slice {
						lists.NewList(steps).Range(func(k int, v interface{}) {
							m := maps.NewMap(v)
							if m.Len() > 0 {
								from := m.GetInt("from")
								to := m.GetInt("to")
								step := m.GetInt("step")
								if from > -1 && to > -1 && step > -1 {
									r := agents.NewScheduleRangeConfig()
									r.From = from
									r.To = to
									r.Step = step
									ranges = append(ranges, r)
								}
							}
						})
					}
				}

				switch timeTypeString {
				case "second":
					schedule.AddSecondRanges(ranges...)
				case "minute":
					schedule.AddMinuteRanges(ranges...)
				case "hour":
					schedule.AddHourRanges(ranges...)
				case "day":
					schedule.AddDayRanges(ranges...)
				case "month":
					schedule.AddMonthRanges(ranges...)
				case "year":
					schedule.AddYearRanges(ranges...)
				case "weekDay":
					schedule.AddWeekDayRanges(ranges...)
				}
			}

			task.AddSchedule(schedule)
		}
	}

	if !params.IsBooting && !params.IsManual && len(task.Schedule) == 0 {
		this.Fail("必须设置一种运行方式：定时、启动或者手动")
	}

	app.AddTask(task)
	err := agent.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 通知更新
	agentutils.PostAgentEvent(agent.Id, agentutils.NewAgentEvent("ADD_TASK", maps.Map{
		"appId":  app.Id,
		"taskId": task.Id,
	}))

	if app.IsSharedWithGroup {
		err := agentutils.SyncApp(agent.Id, agent.GroupIds, app, agentutils.NewAgentEvent("ADD_TASK", maps.Map{
			"appId":  app.Id,
			"taskId": task.Id,
		}), nil)
		if err != nil {
			logs.Error(err)
		}
	}

	this.Success()
}
