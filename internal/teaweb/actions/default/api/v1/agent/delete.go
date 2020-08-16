package agent

import (
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/agentutils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/apiutils"
	"github.com/iwind/TeaGo/actions"
)

type DeleteAction actions.Action

// 删除Agent
func (this *DeleteAction) RunGet(params struct {
	AgentId string
}) {
	if !agentutils.ActionDeleteAgent(params.AgentId, func(message string) {
		apiutils.Fail(this, message)
	}) {
		return
	}

	apiutils.SuccessOK(this)
}
