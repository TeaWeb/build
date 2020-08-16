package waf

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teawaf"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
)

type IndexAction actions.Action

// 策略列表
func (this *IndexAction) RunGet(params struct{}) {
	configs := lists.Map(teaconfigs.SharedWAFList().FindAllConfigs(), func(k int, v interface{}) interface{} {
		config := v.(*teawaf.WAF)
		return maps.Map{
			"id":                  config.Id,
			"name":                config.Name,
			"on":                  config.On,
			"countInboundGroups":  len(config.Inbound),
			"countOutboundGroups": len(config.Outbound),
			"countInboundSets":    config.CountInboundRuleSets(),
			"countOutboundSets":   config.CountOutboundRuleSets(),
		}
	})
	this.Data["configs"] = configs

	this.Show()
}
