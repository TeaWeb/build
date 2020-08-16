package locations

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/locations/locationutils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type WafAction actions.Action

// WAF设置
func (this *WafAction) RunGet(params struct {
	ServerId   string
	LocationId string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	_, location := locationutils.SetCommonInfo(this, params.ServerId, params.LocationId, "waf")

	proxyutils.AddServerMenu(this, true)

	wafList := []maps.Map{}
	for _, waf := range teaconfigs.SharedWAFList().FindAllConfigs() {
		if !waf.On {
			continue
		}
		wafList = append(wafList, maps.Map{
			"id":   waf.Id,
			"name": waf.Name,
		})
	}
	this.Data["wafList"] = wafList

	if len(location.WafId) == 0 {
		this.Data["waf"] = nil
	} else {
		this.Data["waf"] = teaconfigs.SharedWAFList().FindWAF(location.WafId)
		if this.Data["waf"] == nil {
			location.WafId = ""
		}
	}

	this.Data["server"] = proxyutils.WrapServerData(server)
	this.Data["selectedTab"] = "location"

	this.Show()
}
