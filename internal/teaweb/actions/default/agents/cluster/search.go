package cluster

import (
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type SearchAction actions.Action

// 搜索集群
func (this *SearchAction) Run(params struct {
	Rules string
	Port  int
}) {
	if len(params.Rules) == 0 {
		this.FailField("rules", "请输入搜索规则")
	}

	hosts := agentutils.ParseHostRules(params.Rules, 10000)

	result := []maps.Map{}
	for _, host := range hosts {
		result = append(result, maps.Map{
			"addr":  host,
			"state": "WAITING",
		})
	}

	this.Data["hosts"] = result

	this.Success()
}
