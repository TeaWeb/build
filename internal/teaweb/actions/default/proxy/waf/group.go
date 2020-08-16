package waf

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teawaf"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

type GroupAction actions.Action

// 分组信息
func (this *GroupAction) RunGet(params struct {
	WafId   string
	GroupId string
	Inbound bool
}) {
	waf := teaconfigs.SharedWAFList().FindWAF(params.WafId)
	if waf == nil {
		this.Fail("找不到WAF")
	}

	this.Data["config"] = maps.Map{
		"id":            waf.Id,
		"name":          waf.Name,
		"countInbound":  waf.CountInboundRuleSets(),
		"countOutbound": waf.CountOutboundRuleSets(),
	}

	group := waf.FindRuleGroup(params.GroupId)
	if group == nil {
		this.Fail("找不到分组")
	}

	this.Data["inbound"] = group.IsInbound
	this.Data["outbound"] = !group.IsInbound

	this.Data["group"] = group

	// rule sets
	this.Data["sets"] = lists.Map(group.RuleSets, func(k int, v interface{}) interface{} {
		set := v.(*teawaf.RuleSet)

		// 动作说明
		actionLinks := []maps.Map{}
		if set.Action == teawaf.ActionGoGroup {
			nextGroup := waf.FindRuleGroup(set.ActionOptions.GetString("groupId"))
			if nextGroup != nil {
				actionLinks = append(actionLinks, maps.Map{
					"name": nextGroup.Name,
					"url":  "/proxy/waf/group?wafId=" + waf.Id + "&groupId=" + nextGroup.Id,
				})
			}
		} else if set.Action == teawaf.ActionGoSet {
			nextGroup := waf.FindRuleGroup(set.ActionOptions.GetString("groupId"))
			if nextGroup != nil {
				actionLinks = append(actionLinks, maps.Map{
					"name": nextGroup.Name,
					"url":  "/proxy/waf/group?wafId=" + waf.Id + "&groupId=" + nextGroup.Id,
				})

				nextSet := nextGroup.FindRuleSet(set.ActionOptions.GetString("setId"))
				if nextSet != nil {
					actionLinks = append(actionLinks, maps.Map{
						"name": nextSet.Name,
						"url":  "/proxy/waf/group/rule/update?wafId=" + waf.Id + "&groupId=" + nextGroup.Id + "&setId=" + nextSet.Id,
					})
				}
			}
		}

		return maps.Map{
			"id":   set.Id,
			"name": set.Name,
			"rules": lists.Map(set.Rules, func(k int, v interface{}) interface{} {
				rule := v.(*teawaf.Rule)

				return maps.Map{
					"param":             rule.Param,
					"operator":          rule.Operator,
					"value":             rule.Value,
					"isCaseInsensitive": rule.IsCaseInsensitive,
				}
			}),
			"on":            set.On,
			"action":        strings.ToUpper(set.Action),
			"actionOptions": set.ActionOptions,
			"actionName":    teawaf.FindActionName(set.Action),
			"actionLinks":   actionLinks,
			"connector":     strings.ToUpper(set.Connector),
		}
	})

	this.Show()
}
