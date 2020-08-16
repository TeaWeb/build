package app

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/apiutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type AppAction actions.Action

// 单个App信息
func (this *AppAction) RunGet(params struct {
	AgentId string
	AppId   string
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		apiutils.Fail(this, "agent not found")
		return
	}

	app := agent.FindApp(params.AppId)
	if app == nil {
		apiutils.Fail(this, "app not found")
		return
	}

	apiutils.Success(this, maps.Map{
		"config": app,
	})
}
