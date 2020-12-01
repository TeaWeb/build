package waf

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teawaf"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/files"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/rands"
	"github.com/iwind/TeaGo/types"
	"regexp"
)

type AddAction actions.Action

// 添加策略
func (this *AddAction) RunGet(params struct{}) {
	this.Data["groups"] = lists.Map(teawaf.Template().Inbound, func(k int, v interface{}) interface{} {
		g := v.(*teawaf.RuleGroup)
		return maps.Map{
			"name":      g.Name,
			"code":      g.Code,
			"isChecked": g.On,
		}
	})

	// 匹配条件运算符
	this.Data["condOperators"] = shared.AllRequestOperators()
	this.Data["condVariables"] = proxyutils.DefaultRequestVariables()

	this.Show()
}

// 保存提交
func (this *AddAction) RunPost(params struct {
	Name       string
	GroupCodes []string

	On bool

	BlockStatusCode string
	BlockBody       string
	BlockURL        string

	Must *actions.Must
}) {
	params.Must.
		Field("name", params.Name).
		Require("请输入策略名称")

	if len(params.BlockStatusCode) > 0 && !regexp.MustCompile(`^\d{3}$`).MatchString(params.BlockStatusCode) {
		this.FailField("blockStatusCode", "请输入正确的HTTP状态码")
	}
	statusCode := types.Int(params.BlockStatusCode)

	waf := teawaf.NewWAF()
	waf.Name = params.Name
	waf.On = params.On
	waf.ActionBlock = &teawaf.BlockAction{
		StatusCode: statusCode,
		Body:       params.BlockBody,
		URL:        params.BlockURL,
	}

	// 匹配条件
	conds, breakCond, err := proxyutils.ParseRequestConds(this.Request, "request")
	if err != nil {
		this.Fail("匹配条件\"" + breakCond.Param + " " + breakCond.Operator + " " + breakCond.Value + "\"校验失败：" + err.Error())
	}
	waf.Cond = conds

	template := teawaf.Template()

	for _, g := range template.Inbound {
		newGroup := teawaf.NewRuleGroup()
		newGroup.Id = rands.HexString(16)
		newGroup.On = lists.ContainsString(params.GroupCodes, g.Code)
		newGroup.Code = g.Code
		newGroup.Name = g.Name
		newGroup.RuleSets = g.RuleSets
		newGroup.IsInbound = g.IsInbound
		newGroup.Description = g.Description
		waf.AddRuleGroup(newGroup)
	}

	filename := "waf." + waf.Id + ".conf"
	err = waf.Save(Tea.ConfigFile(filename))
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	wafList := teaconfigs.SharedWAFList()
	wafList.AddFile(filename)
	err = wafList.Save()
	if err != nil {
		err1 := files.NewFile(Tea.ConfigFile(filename)).Delete()
		if err1 != nil {
			logs.Error(err1)
		}

		this.Fail("保存失败：" + err.Error())
	}

	this.Success()
}
