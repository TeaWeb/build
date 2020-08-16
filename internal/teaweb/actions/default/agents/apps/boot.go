package apps

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
)

type BootAction actions.Action

// 启动任务
func (this *BootAction) Run(params struct {
	AgentId string
	AppId   string
}) {
	app := agentutils.InitAppData(this, params.AgentId, params.AppId, "boot")
	this.Data["tasks"] = lists.Map(app.FindBootingTasks(), func(k int, v interface{}) interface{} {
		return agentutils.FormatTask(v.(*agents.TaskConfig), params.AgentId)
	})

	this.Show()
}
