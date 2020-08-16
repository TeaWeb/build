package groups

import (
	"github.com/TeaWeb/build/internal/teaconfigs/agents"
	"github.com/TeaWeb/build/internal/teaconfigs/notices"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

type IndexAction actions.Action

// 分组管理
func (this *IndexAction) Run(params struct{}) {
	groups := []maps.Map{}
	for _, group := range agents.SharedGroupList().Groups {
		countReceivers := 0
		for _, receivers := range group.NoticeSetting {
			countReceivers += len(receivers)
		}

		groups = append(groups, maps.Map{
			"id":             group.Id,
			"name":           group.Name,
			"on":             group.On,
			"countAgents":    group.CountAgents,
			"countReceivers": countReceivers,
			"canDelete":      !group.IsDefault,
			"maxAgents":      group.MaxAgents,
			"dayFrom":        group.DayFrom,
			"dayTo":          group.DayTo,
			"isExpired":      len(group.DayTo) > 0 && group.DayTo < timeutil.Format("Y-m-d"),
		})
	}

	this.Data["groups"] = groups
	this.Data["noticeLevels"] = notices.AllNoticeLevels()

	this.Show()
}
