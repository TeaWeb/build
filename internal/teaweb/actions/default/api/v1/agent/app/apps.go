package app

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/apiutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type AppsAction actions.Action

// 某个Agent下所有App
func (this *AppsAction) RunGet(params struct {
	AgentId string
}) {
	agent := agents.NewAgentConfigFromId(params.AgentId)
	if agent == nil {
		apiutils.Fail(this, "agent not found")
		return
	}

	result := []maps.Map{}
	for _, app := range agent.Apps {
		result = append(result, maps.Map{
			"config": app,
		})
	}

	apiutils.Success(this, result)
}
