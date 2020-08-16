package waf

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
)

type GroupUpdateAction actions.Action

// 修改分组
func (this *GroupUpdateAction) RunGet(params struct {
	WafId   string
	GroupId string
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

	this.Show()
}

// 保存修改
func (this *GroupUpdateAction) RunPost(params struct {
	WafId       string
	GroupId     string
	Name        string
	Description string
	On          bool
	Must        *actions.Must
}) {
	wafList := teaconfigs.SharedWAFList()
	config := wafList.FindWAF(params.WafId)
	if config == nil {
		this.Fail("找不到WAF")
	}

	params.Must.
		Field("name", params.Name).
		Require("请输入分组名称")

	group := config.FindRuleGroup(params.GroupId)
	if group == nil {
		this.Fail("找不到分组")
	}

	group.On = params.On
	group.Name = params.Name
	group.Description = params.Description

	err := wafList.SaveWAF(config)
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
