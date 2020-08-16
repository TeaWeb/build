package apps

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
)

type ManualAction actions.Action

// 手动任务
func (this *ManualAction) Run(params struct {
	AgentId string
	AppId   string
}) {
	app := agentutils.InitAppData(this, params.AgentId, params.AppId, "manual")
	this.Data["tasks"] = lists.Map(app.FindManualTasks(), func(k int, v interface{}) interface{} {
		return agentutils.FormatTask(v.(*agents.TaskConfig), params.AgentId)
	})

	this.Show()
}
