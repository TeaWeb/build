package apps

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/TeaWeb/build/internal/teadb"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction actions.Action

// 看板首页
func (this *IndexAction) Run(params struct {
	AgentId string
}) {
	this.Data["agentId"] = params.AgentId

	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		this.Fail("找不到要修改的Agent")
	}

	// 用户自定义App
	this.Data["apps"] = lists.Map(agent.Apps, func(k int, v interface{}) interface{} {
		app := v.(*agents.AppConfig)

		// 最新一条数据
		level := notices.NoticeLevelNone
		for _, item := range app.Items {
			if !item.On {
				continue
			}
			value, err := teadb.AgentValueDAO().FindLatestItemValue(agent.Id, app.Id, item.Id)
			if err == nil && value != nil {
				if value.NoticeLevel == notices.NoticeLevelWarning || value.NoticeLevel == notices.NoticeLevelError && value.NoticeLevel > level {
					level = value.NoticeLevel
				}
			}
		}

		return maps.Map{
			"on":                app.On,
			"id":                app.Id,
			"name":              app.Name,
			"items":             app.Items,
			"bootingTasks":      app.FindBootingTasks(),
			"manualTasks":       app.FindManualTasks(),
			"schedulingTasks":   app.FindSchedulingTasks(),
			"isSharedWithGroup": app.IsSharedWithGroup,
			"isWarning":         level == notices.NoticeLevelWarning,
			"isError":           level == notices.NoticeLevelError,
		}
	})

	this.Show()
}
