package waf

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teawaf"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/waf/wafutils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/rands"
	"github.com/iwind/TeaGo/types"
	"net/http"
	"regexp"
)

type UpdateAction actions.Action

// 修改
func (this *UpdateAction) RunGet(params struct {
	WafId string
}) {
	waf := teaconfigs.SharedWAFList().FindWAF(params.WafId)
	if waf == nil {
		this.Fail("找不到WAF")
	}

	if waf.ActionBlock == nil {
		waf.ActionBlock = &teawaf.BlockAction{
			StatusCode: http.StatusForbidden,
		}
	}

	this.Data["config"] = maps.Map{
		"id":            waf.Id,
		"name":          waf.Name,
		"on":            waf.On,
		"countInbound":  waf.CountInboundRuleSets(),
		"countOutbound": waf.CountOutboundRuleSets(),
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

	// 匹配条件运算符
	this.Data["condOperators"] = shared.AllRequestOperators()
	this.Data["condVariables"] = proxyutils.DefaultRequestVariables()

	this.Show()
}

// 保存修改
func (this *UpdateAction) RunPost(params struct {
	WafId      string
	Name       string
	GroupCodes []string

	On bool

	BlockStatusCode string
	BlockBody       string
	BlockURL        string

	Must *actions.Must
}) {
	waf := teaconfigs.SharedWAFList().FindWAF(params.WafId)
	if waf == nil {
		this.Fail("无法找到WAF")
	}

	params.Must.
		Field("name", params.Name).
		Require("请输入策略名称")

	if len(params.BlockStatusCode) > 0 && !regexp.MustCompile(`^\d{3}$`).MatchString(params.BlockStatusCode) {
		this.FailField("blockStatusCode", "请输入正确的HTTP状态码")
	}
	statusCode := types.Int(params.BlockStatusCode)

	waf.Name = params.Name
	waf.On = params.On
	waf.ActionBlock = &teawaf.BlockAction{
		StatusCode: statusCode,
		Body:       params.BlockBody,
		URL:        params.BlockURL,
	}

	// add new group
	template := teawaf.Template()
	for _, groupCode := range params.GroupCodes {
		g := waf.FindRuleGroupWithCode(groupCode)
		if g != nil {
			g.On = true
			continue
		}
		g = template.FindRuleGroupWithCode(groupCode)
		g.Id = rands.HexString(16)
		g.On = true
		waf.AddRuleGroup(g)
	}

	// remove old group {
	for _, g := range waf.Inbound {
		if len(g.Code) > 0 && !lists.ContainsString(params.GroupCodes, g.Code) {
			g.On = false
			continue
		}
	}

	for _, g := range waf.Outbound {
		if len(g.Code) > 0 && !lists.ContainsString(params.GroupCodes, g.Code) {
			g.On = false
			continue
		}
	}

	// 匹配条件
	conds, breakCond, err := proxyutils.ParseRequestConds(this.Request, "request")
	if err != nil {
		this.Fail("匹配条件\"" + breakCond.Param + " " + breakCond.Operator + " " + breakCond.Value + "\"校验失败：" + err.Error())
	}
	waf.Cond = conds

	filename := "waf." + waf.Id + ".conf"
	err = waf.Save(Tea.ConfigFile(filename))
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	// 通知刷新
	if wafutils.IsPolicyUsed(waf.Id) {
		proxyutils.NotifyChange()
	}

	this.Success()
}
