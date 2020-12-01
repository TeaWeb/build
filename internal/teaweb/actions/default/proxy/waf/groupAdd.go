package waf

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teawaf"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/rands"
)

type GroupAddAction actions.Action

// 添加分组
func (this *GroupAddAction) RunGet(params struct {
	WafId   string
	Inbound bool
}) {
	config := teaconfigs.SharedWAFList().FindWAF(params.WafId)
	if config == nil {
		this.Fail("找不到WAF")
	}

	this.Data["inbound"] = params.Inbound
	this.Data["outbound"] = !params.Inbound
	this.Data["config"] = maps.Map{
		"id":            config.Id,
		"name":          config.Name,
		"countInbound":  config.CountInboundRuleSets(),
		"countOutbound": config.CountOutboundRuleSets(),
	}

	this.Show()
}

// 保存分组
func (this *GroupAddAction) RunPost(params struct {
	WafId       string
	Name        string
	Description string
	On          bool
	Inbound     bool
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

	group := teawaf.NewRuleGroup()
	group.Id = rands.HexString(16)
	group.On = params.On
	group.Name = params.Name
	group.Description = params.Description
	group.IsInbound = params.Inbound

	config.AddRuleGroup(group)
	err := wafList.SaveWAF(config)
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	this.Data["id"] = group.Id
	this.Success()
}
