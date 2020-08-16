package agentutils

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/utils/time"
	"time"
)

// App菜单
func InitAppData(actionWrapper actions.ActionWrapper, agentId string, appId string, tabbar string) *agents.AppConfig {
	agent := agents.NewAgentConfigFromId(agentId)
	action := actionWrapper.Object()
	if agent == nil {
		action.Fail("找不到Agent")
	}

	app := agent.FindApp(appId)
	if app == nil {
		action.Fail("找不到App")
	}

	action.Data["agentId"] = agentId
	action.Data["app"] = maps.Map{
		"id":                   app.Id,
		"name":                 app.Name,
		"on":                   app.On,
		"countItems":           len(app.Items),
		"countBootTasks":       len(app.FindBootingTasks()),
		"countScheduleTasks":   len(app.FindSchedulingTasks()),
		"countManualTasks":     len(app.FindManualTasks()),
		"countNoticeReceivers": app.CountNoticeReceivers(),
		"isSharedWithGroup":    app.IsSharedWithGroup,
	}
	action.Data["selectedTabbar"] = tabbar

	return app
}

// 格式化任务信息
func FormatTask(task *agents.TaskConfig, agentId string) maps.Map {
	// 最近执行
	processLog, err := teadb.AgentLogDAO().FindLatestTaskLog(agentId, task.Id)
	runTime := ""
	if err != nil {
		logs.Error(err)
	} else {
		if processLog != nil {
			runTime = timeutil.Format("Y-m-d H:i:s", time.Unix(processLog.Timestamp, 0))
		}
	}

	return maps.Map{
		"id":        task.Id,
		"on":        task.On,
		"name":      task.Name,
		"script":    task.Script,
		"isBooting": task.IsBooting,
		"isManual":  task.IsManual,
		"schedules": lists.Map(task.Schedule, func(k int, v interface{}) interface{} {
			return v.(*agents.ScheduleConfig).Summary()
		}),
		"runTime": runTime,
	}
}

// 查找共享的Agent
func FindSharedAgents(currentAgentId string, groupIds []string, app *agents.AppConfig) []*agents.AgentConfig {
	result := []*agents.AgentConfig{}
	if app.IsSharedWithGroup {
		for _, agent := range agents.AllSharedAgents() {
			if agent.Id == currentAgentId || !agent.BelongsToGroups(groupIds) {
				continue
			}
			if !agent.HasApp(app.Id) {
				continue
			}
			result = append(result, agent)
		}
	}
	return result

}

// 同步App到其他Agents
// op是附加操作
func SyncApp(currentAgentId string, groupIds []string, app *agents.AppConfig, event *AgentEvent, op func(agent *agents.AgentConfig) error) error {
	if app.IsSharedWithGroup { // 添加共享
		for _, agent := range agents.AllSharedAgents() {
			if agent.Id == currentAgentId || !agent.BelongsToGroups(groupIds) {
				continue
			}
			agent.ReplaceApp(app)
			if op != nil {
				err := op(agent)
				if err != nil {
					return err
				}
			}
			err := agent.Save()
			if err != nil {
				return err
			}
			if event != nil {
				PostAgentEvent(agent.Id, event)
			}
		}
	} else { // 取消共享，需要删除其他Agent中的App
		for _, agent := range agents.AllSharedAgents() {
			if agent.Id == currentAgentId {
				continue
			}
			if !agent.HasApp(app.Id) {
				continue
			}

			if op != nil {
				err := op(agent)
				if err != nil {
					return err
				}
			}

			// 删除图表
			board := agents.NewAgentBoard(agent.Id)
			if board != nil {
				board.RemoveApp(app.Id)
				err := board.Save()
				if err != nil {
					return err
				}
			}

			// 删除App
			agent.RemoveApp(app.Id)

			err := agent.Save()
			if err != nil {
				return err
			}
			if event != nil {
				PostAgentEvent(agent.Id, event)
			}
		}
	}

	return nil
}

// 仅同步Event
func SyncAppEvent(currentAgentId string, groupIds []string, app *agents.AppConfig, event *AgentEvent) error {
	if app.IsSharedWithGroup {
		for _, agent := range agents.AllSharedAgents() {
			if agent.Id == currentAgentId || !agent.BelongsToGroups(groupIds) {
				continue
			}
			if !agent.HasApp(app.Id) {
				continue
			}
			PostAgentEvent(agent.Id, event)
		}
	}

	return nil
}

// 仅同步Chart
func SyncRemoveChart(currentAgentId string, groupIds []string, app *agents.AppConfig, chartId string) error {
	if app.IsSharedWithGroup {
		for _, agent := range agents.AllSharedAgents() {
			if agent.Id == currentAgentId || !agent.BelongsToGroups(groupIds) {
				continue
			}
			if !agent.HasApp(app.Id) {
				continue
			}
			board := agents.NewAgentBoard(agent.Id)
			board.RemoveChart(chartId)
			err := board.Save()
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 添加Chart
func SyncAddChart(currentAgentId string, groupIds []string, app *agents.AppConfig, itemId string, chartId string) error {
	if app.IsSharedWithGroup {
		for _, agent := range agents.AllSharedAgents() {
			if agent.Id == currentAgentId || !agent.BelongsToGroups(groupIds) {
				continue
			}
			if !agent.HasApp(app.Id) {
				continue
			}
			board := agents.NewAgentBoard(agent.Id)
			board.AddChart(app.Id, itemId, chartId)
			err := board.Save()
			if err != nil {
				return err
			}
		}
	}
	return nil
}
