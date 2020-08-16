package servers

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
)

type WafUpdateAction actions.Action

// WAF修改
func (this *WafUpdateAction) RunPost(params struct {
	ServerId string
	WafId    string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	if len(params.WafId) == 0 {
		server.WAFOn = true // always true
		server.WafId = ""
	} else {
		waf := teaconfigs.SharedWAFList().FindWAF(params.WafId)
		if waf == nil {
			this.Fail("要设置的WAF不存在")
		}

		server.WAFOn = true
		server.WafId = waf.Id
	}

	err := server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	proxyutils.NotifyChange()

	this.Success()
}
