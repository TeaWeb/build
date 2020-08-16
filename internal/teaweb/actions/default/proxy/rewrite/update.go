package rewrite

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/locations/locationutils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/Tea"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/maps"
	"regexp"
	"strings"
)

type UpdateAction actions.Action

// 修改重写规则
func (this *UpdateAction) Run(params struct {
	From       string
	ServerId   string
	LocationId string
	RewriteId  string
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	this.Data["server"] = proxyutils.WrapServerData(server)
	this.Data["from"] = params.From
	this.Data["locationId"] = params.LocationId

	this.Data["location"] = nil

	if len(params.LocationId) > 0 {
		location := server.FindLocation(params.LocationId)
		if location != nil {
			locationutils.SetCommonInfo(this, params.ServerId, params.LocationId, "rewrite")
		}
	}

	// 已经有的代理服务
	proxyConfigs := teaconfigs.LoadServerConfigsFromDir(Tea.ConfigDir())
	proxies := []maps.Map{}
	for _, proxyConfig := range proxyConfigs {
		if proxyConfig.Id == server.Id {
			continue
		}

		name := proxyConfig.Description
		if !proxyConfig.On {
			name += "(未启用)"
		}
		proxies = append(proxies, maps.Map{
			"id":   proxyConfig.Id,
			"name": name,
		})
	}
	this.Data["proxies"] = proxies

	this.Data["typeOptions"] = []maps.Map{
		{
			"name":  "匹配前缀",
			"value": teaconfigs.LocationPatternTypePrefix,
		},
		{
			"name":  "精准匹配",
			"value": teaconfigs.LocationPatternTypeExact,
		},
		{
			"name":  "正则表达式匹配",
			"value": teaconfigs.LocationPatternTypeRegexp,
		},
	}

	// 运算符
	this.Data["operators"] = shared.AllRequestOperators()

	// 当前Rewrite信息
	rewriteList, err := server.FindRewriteList(params.LocationId)
	if err != nil {
		this.Fail(err.Error())
	}
	rewrite := rewriteList.FindRewriteRule(params.RewriteId)
	if rewrite == nil {
		this.Fail("找不到要修改的重写规则")
	}
	this.Data["rewrite"] = maps.Map{
		"on":           rewrite.On,
		"id":           rewrite.Id,
		"pattern":      rewrite.Pattern,
		"replace":      rewrite.TargetURL(),
		"flags":        rewrite.Flags,
		"proxyId":      rewrite.TargetProxy(),
		"conds":        rewrite.Cond,
		"targetType":   rewrite.TargetType(),
		"redirectMode": rewrite.RedirectMode(),
		"isBreak":      rewrite.IsBreak,
		"isPermanent":  rewrite.IsPermanent,
		"proxyHost":    rewrite.ProxyHost,
	}

	// 变量
	this.Data["variables"] = proxyutils.DefaultRequestVariables()

	this.Show()
}

// 提交保存
func (this *UpdateAction) RunPost(params struct {
	ServerId     string
	LocationId   string
	RewriteId    string
	On           bool
	Pattern      string
	Replace      string
	ProxyId      string
	TargetType   string
	RedirectMode string
	IsBreak      bool
	IsPermanent  bool
	ProxyHost    string
	Must         *actions.Must
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	params.Must.
		Field("pattern", params.Pattern).
		Require("请输入匹配规则").
		Expect(func() (message string, success bool) {
			_, err := regexp.Compile(params.Pattern)
			if err != nil {
				return "匹配规则错误：" + err.Error(), false
			}
			return "", true
		}).

		Field("targetType", params.TargetType).
		In([]string{"url", "proxy"}, "目标类型错误")

	if params.TargetType == "proxy" {
		params.Must.
			Field("proxyId", params.ProxyId).
			Require("请选择目标代理")
	}
	params.Must.
		Field("replace", params.Replace).
		Require("请输入目标URL")

	if len(params.Replace) == 0 {
		params.Replace = "/"
	}

	rewriteList, err := server.FindRewriteList(params.LocationId)
	if err != nil {
		this.Fail(err.Error())
	}

	rewriteRule := rewriteList.FindRewriteRule(params.RewriteId)
	if rewriteRule == nil {
		this.Fail("找不到要修改的Rewrite")
	}

	rewriteRule.On = params.On
	rewriteRule.Pattern = params.Pattern
	if params.TargetType == "url" {
		rewriteRule.Replace = params.Replace
	} else {
		rewriteRule.Replace = "proxy://" + params.ProxyId + "/" + strings.TrimLeft(params.Replace, "/")
	}
	rewriteRule.Flags = []string{}
	rewriteRule.FlagOptions = maps.Map{}
	if len(params.RedirectMode) > 0 {
		rewriteRule.AddFlag(params.RedirectMode, nil)
	}

	// 匹配条件
	conds, breakCond, err := proxyutils.ParseRequestConds(this.Request, "request")
	if err != nil {
		this.Fail("匹配条件\"" + breakCond.Param + " " + breakCond.Operator + " " + breakCond.Value + "\"校验失败：" + err.Error())
	}
	rewriteRule.Cond = conds

	rewriteRule.IsBreak = params.IsBreak
	rewriteRule.IsPermanent = params.IsPermanent
	rewriteRule.ProxyHost = params.ProxyHost

	err = server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	proxyutils.NotifyChange()

	this.Success()
}
