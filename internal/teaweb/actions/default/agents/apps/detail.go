package apps

import (
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
)

type DetailAction actions.Action

// App详情
func (this *DetailAction) Run(params struct {
	AgentId string
	AppId   string
}) {
	agentutils.InitAppData(this, params.AgentId, params.AppId, "detail")

	this.Show()
}
