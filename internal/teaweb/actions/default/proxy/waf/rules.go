package waf

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teawaf"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

type RulesAction actions.Action

// 规则
func (this *RulesAction) RunGet(params struct {
	Inbound bool
	WafId   string
}) {
	config := teaconfigs.SharedWAFList().FindWAF(params.WafId)
	if config == nil {
		this.Fail("找不到WAF")
	}

	this.Data["config"] = maps.Map{
		"id":            config.Id,
		"name":          config.Name,
		"countInbound":  config.CountInboundRuleSets(),
		"countOutbound": config.CountOutboundRuleSets(),
	}
	this.Data["inbound"] = params.Inbound
	this.Data["outbound"] = !params.Inbound

	if params.Inbound {
		this.Data["groups"] = lists.Map(config.Inbound, func(k int, v interface{}) interface{} {
			group := v.(*teawaf.RuleGroup)
			return maps.Map{
				"id":            group.Id,
				"code":          group.Code,
				"name":          group.Name,
				"description":   group.Description,
				"on":            group.On,
				"countRuleSets": len(group.RuleSets),
				"canDelete":     len(group.Code) == 0,
			}
		})
	} else {
		this.Data["groups"] = lists.Map(config.Outbound, func(k int, v interface{}) interface{} {
			group := v.(*teawaf.RuleGroup)
			return maps.Map{
				"id":            group.Id,
				"code":          group.Code,
				"name":          group.Name,
				"description":   group.Description,
				"on":            group.On,
				"countRuleSets": len(group.RuleSets),
				"canDelete":     len(group.Code) == 0,
			}
		})
	}

	this.Show()
}
