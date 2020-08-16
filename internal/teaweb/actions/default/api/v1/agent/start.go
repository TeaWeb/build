package agent

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/apiutils"
	"github.com/iwind/TeaGo/actions"
)

type StartAction actions.Action

// 启动Agent
func (this *StartAction) RunGet(params struct {
	AgentId string
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		apiutils.Fail(this, "not found")
		return
	}

	agent.On = true
	err := agent.Save()
	if err != nil {
		apiutils.Fail(this, "保存失败:"+err.Error())
		return
	}

	apiutils.SuccessOK(this)
}
