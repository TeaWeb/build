package waf

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teawaf"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"gopkg.in/yaml.v3"
	"strings"
)

type ExportAction actions.Action

// 导出
func (this *ExportAction) RunGet(params struct {
	WafId string

	Export   bool
	GroupIds string
}) {
	waf := teaconfigs.SharedWAFList().FindWAF(params.WafId)
	if waf == nil {
		this.Fail("找不到WAF")
	}

	// 导出
	if params.Export {
		waf1 := waf.Copy()
		waf1.Inbound = []*teawaf.RuleGroup{}
		waf1.Outbound = []*teawaf.RuleGroup{}
		if len(params.GroupIds) > 0 {
			groupIds := strings.Split(params.GroupIds, ",")
			for _, groupId := range groupIds {
				group := waf.FindRuleGroup(groupId)
				if group == nil {
					continue
				}
				waf1.AddRuleGroup(group)
			}
		}

		data, err := yaml.Marshal(waf1)
		if err != nil {
			this.WriteString(err.Error())
			return
		}

		filename := "waf." + waf1.Id + ".conf"
		this.ResponseWriter.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")

		this.Write(data)

		return
	}

	this.Data["config"] = maps.Map{
		"id":            waf.Id,
		"name":          waf.Name,
		"countInbound":  waf.CountInboundRuleSets(),
		"countOutbound": waf.CountOutboundRuleSets(),
	}

	groups := []maps.Map{}
	for _, group := range waf.Inbound {
		groups = append(groups, maps.Map{
			"id":        group.Id,
			"name":      "[入站]" + group.Name,
			"countSets": len(group.RuleSets),
			"on":        group.On,
		})
	}
	for _, group := range waf.Outbound {
		groups = append(groups, maps.Map{
			"id":        group.Id,
			"name":      "[出站]" + group.Name,
			"countSets": len(group.RuleSets),
			"on":        group.On,
		})
	}

	this.Data["groups"] = groups

	this.Show()
}
