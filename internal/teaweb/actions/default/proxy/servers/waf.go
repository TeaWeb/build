package servers

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type WafAction actions.Action

// WAF设置
func (this *WafAction) RunGet(params struct {
	ServerId string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

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

	if len(server.WafId) == 0 {
		this.Data["waf"] = nil
	} else {
		this.Data["waf"] = teaconfigs.SharedWAFList().FindWAF(server.WafId)
		if this.Data["waf"] == nil {
			server.WafId = ""
		}
	}

	this.Data["server"] = proxyutils.WrapServerData(server)
	this.Data["selectedTab"] = "waf"

	this.Show()
}
