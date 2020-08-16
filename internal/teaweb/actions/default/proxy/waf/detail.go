package waf

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teawaf"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

type DetailAction actions.Action

// 详情
func (this *DetailAction) RunGet(params struct {
	WafId string
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
		"on":            waf.On,
		"actionBlock":   waf.ActionBlock,
		"cond":          waf.Cond,
	}

	this.Data["groups"] = lists.Map(teawaf.Template().Inbound, func(k int, v interface{}) interface{} {
		g := v.(*teawaf.RuleGroup)
		group := waf.FindRuleGroupWithCode(g.Code)

		return maps.Map{
			"name":      g.Name,
			"code":      g.Code,
			"isChecked": group != nil && group.On,
		}
	})

	// 正在使用此策略的项目
	configItems := []maps.Map{}
	serverList, _ := teaconfigs.SharedServerList()
	if serverList != nil {
		for _, server := range serverList.FindAllServers() {

			if server.WafId == waf.Id {
				configItems = append(configItems, maps.Map{
					"type":   "server",
					"server": server.Description,
					"link":   "/proxy/servers/waf?serverId=" + server.Id,
				})
			}

			for _, location := range server.Locations {
				if location.WafId == waf.Id {
					configItems = append(configItems, maps.Map{
						"type":     "location",
						"server":   server.Description,
						"location": location.Pattern,
						"link":     "/proxy/locations/waf?serverId=" + server.Id + "&locationId=" + location.Id,
					})
				}
			}
		}
	}

	this.Data["configItems"] = configItems

	// 是否有新的模版变更
	this.Data["newItems"] = waf.MergeTemplate()

	this.Show()
}
