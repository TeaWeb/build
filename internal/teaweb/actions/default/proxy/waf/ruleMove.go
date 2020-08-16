package waf

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/waf/wafutils"
	"github.com/iwind/TeaGo/actions"
)

type RuleMoveAction actions.Action

// 移动规则集
func (this *RuleMoveAction) RunPost(params struct {
	WafId     string
	GroupId   string
	FromIndex int
	ToIndex   int
}) {
	wafList := teaconfigs.SharedWAFList()
	waf := wafList.FindWAF(params.WafId)
	if waf == nil {
		this.Fail("找不到WAF")
	}

	group := waf.FindRuleGroup(params.GroupId)
	if group == nil {
		this.Fail("找不到分组")
	}

	group.MoveRuleSet(params.FromIndex, params.ToIndex)
	err := wafList.SaveWAF(waf)
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 通知刷新
	if wafutils.IsPolicyUsed(waf.Id) {
		proxyutils.NotifyChange()
	}

	this.Success()
}
