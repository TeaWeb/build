package groups

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/agents/agentutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type DetailAction actions.Action

// 详情
func (this *DetailAction) Run(params struct {
	GroupId string
}) {
	group := agents.SharedGroupList().FindGroup(params.GroupId)
	if group == nil {
		this.Fail("找不到Group")
	}
	this.Data["group"] = group
	this.Data["isExpired"] = len(group.DayTo) > 0 && group.DayTo < timeutil.Format("Y-m-d")

	// Agents列表
	groupAgents := []maps.Map{}
	agentList, err := agents.SharedAgentList()
	groupKeys := []maps.Map{}
	if err != nil {
		logs.Error(err)
	} else {
		allAgents := agentList.FindAllAgents()
		for _, a := range allAgents {
			if a.BelongsToGroup(params.GroupId) {
				state := agentutils.FindAgentState(a.Id)
				groupAgents = append(groupAgents, maps.Map{
					"on":        a.On,
					"id":        a.Id,
					"name":      a.Name,
					"host":      a.Host,
					"isWaiting": state.IsActive,
				})
			}
		}

		for _, key := range group.Keys {
			groupKeys = append(groupKeys, maps.Map{
				"key":         key.Key,
				"dayFrom":     key.DayFrom,
				"dayTo":       key.DayTo,
				"maxAgents":   key.MaxAgents,
				"countAgents": key.CountAgents,
			})
		}
	}

	this.Data["groupKeys"] = groupKeys
	this.Data["agents"] = groupAgents

	this.Show()
}
