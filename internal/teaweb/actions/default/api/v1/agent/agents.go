package agent

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/api/apiutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type AgentsAction actions.Action

// Agent列表
func (this *AgentsAction) RunGet(params struct{}) {
	result := []maps.Map{}
	for _, agent := range agents.SharedAgents() {
		result = append(result, maps.Map{
			"config": agent,
		})
	}
	apiutils.Success(this, result)
}
