package waf

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/waf/wafutils"
	"github.com/iwind/TeaGo/actions"
)

type GroupMoveAction actions.Action

// 移动分组
func (this *GroupMoveAction) RunPost(params struct {
	WafId     string
	Inbound   bool
	FromIndex int
	ToIndex   int
}) {
	wafList := teaconfigs.SharedWAFList()
	waf := wafList.FindWAF(params.WafId)
	if waf == nil {
		this.Fail("找不到WAF")
	}

	if params.Inbound {
		waf.MoveInboundRuleGroup(params.FromIndex, params.ToIndex)
		err := wafList.SaveWAF(waf)
		if err != nil {
			this.Fail("保存失败：" + err.Error())
		}
	} else {
		waf.MoveOutboundRuleGroup(params.FromIndex, params.ToIndex)
		err := wafList.SaveWAF(waf)
		if err != nil {
			this.Fail("保存失败：" + err.Error())
		}
	}

	// 通知刷新
	if wafutils.IsPolicyUsed(waf.Id) {
		proxyutils.NotifyChange()
	}

	this.Success()
}
