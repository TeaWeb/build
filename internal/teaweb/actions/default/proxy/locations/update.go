package locations

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/locations/locationutils"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/TeaWeb/build/internal/teaweb/configs"
	"github.com/iwind/TeaGo/actions"
	"github.com/iwind/TeaGo/lists"
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"regexp"
	"strconv"
	"strings"
)

type UpdateAction actions.Action

// 修改
func (this *UpdateAction) Run(params struct {
	ServerId    string
	LocationId  string
	From        string
	ShowSpecial bool
}) {
	_, location := locationutils.SetCommonInfo(this, params.ServerId, params.LocationId, "detail")

	this.Data["from"] = params.From
	this.Data["showSpecial"] = params.ShowSpecial

	this.Data["patternTypes"] = teaconfigs.AllLocationPatternTypes()
	this.Data["usualCharsets"] = teautils.UsualCharsets
	this.Data["charsets"] = teautils.AllCharsets
	this.Data["accessLogIsInherited"] = len(location.AccessLog) == 0
	if len(location.AccessLog) == 0 {
		location.AccessLog = []*teaconfigs.AccessLogConfig{teaconfigs.NewAccessLogConfig()}
	}
	this.Data["accessLogs"] = proxyutils.FormatAccessLog(location.AccessLog)

	if len(location.Pages) == 0 {
		location.Pages = []*teaconfigs.PageConfig{}
	}

	//gzip
	if location.Gzip == nil {
		this.Data["gzip"] = &teaconfigs.GzipConfig{
			Level:     -1,
			MinLength: "",
			MimeTypes: teaconfigs.DefaultGzipMimeTypes,
		}
	} else {
		this.Data["gzip"] = location.Gzip
	}

	this.Data["location"] = maps.Map{
		"id":                location.Id,
		"on":                location.On,
		"pattern":           location.PatternString(),
		"type":              location.PatternType(),
		"name":              location.Name,
		"isBreak":           location.IsBreak,
		"isReverse":         location.IsReverse(),
		"isCaseInsensitive": location.IsCaseInsensitive(),
		"root":              location.Root,
		"urlPrefix":         location.URLPrefix,
		"index":             location.Index,
		"charset":           location.Charset,
		"maxBodySize":       location.MaxBodySize,
		"enableStat":        !location.DisableStat,
		"redirectToHttps":   location.RedirectToHttps,
		"conds":             location.Cond,
		"denyConds":         location.DenyCond,
		"denyAll":           location.DenyAll,

		// 菜单用
		"rewrite":     location.Rewrite,
		"headers":     location.Headers,
		"fastcgi":     location.Fastcgi,
		"cachePolicy": location.CachePolicy,
		"websocket":   location.Websocket,
		"backends":    location.Backends,
		"wafOn":       location.WAFOn,
		"wafId":       location.WafId,

		"shutdown": location.Shutdown,
		"pages":    location.Pages,
	}

	// 运算符
	this.Data["operators"] = shared.AllRequestOperators()

	// 变量
	this.Data["variables"] = proxyutils.DefaultRequestVariables()

	// 目录补全
	security := configs.SharedAdminConfig().Security
	if security != nil {
		this.Data["dirAutoComplete"] = security.DirAutoComplete
	} else {
		this.Data["dirAutoComplete"] = false
	}

	this.Show()
}

// 保存修改
func (this *UpdateAction) RunPost(params struct {
	ServerId    string
	LocationId  string
	Pattern     string
	PatternType int

	IsBreak bool

	Name                 string
	Root                 string
	URLPrefix            string `alias:"urlPrefix"`
	Charset              string
	Index                []string
	MaxBodySize          float64
	MaxBodyUnit          string
	AccessLogIsInherited bool
	EnableStat           bool
	DenyAll              bool

	GzipLevel          int8
	GzipMinLength      float64
	GzipMinUnit        string
	GzipMimeTypeValues []string

	RedirectToHttps   bool
	On                bool
	IsReverse         bool
	IsCaseInsensitive bool

	PageStatusList    []string
	PageURLList       []string
	PageNewStatusList []string

	ShutdownPageOn     bool
	ShutdownPageURL    string
	ShutdownPageStatus int
}) {
	server := teaconfigs.NewServerConfigFromId(params.ServerId)
	if server == nil {
		this.Fail("找不到Server")
	}

	location := server.FindLocation(params.LocationId)
	if location == nil {
		this.Fail("找不到要修改的Location")
	}

	// 校验正则
	if params.PatternType == teaconfigs.LocationPatternTypeRegexp {
		_, err := regexp.Compile(params.Pattern)
		if err != nil {
			this.Fail("正则表达式校验失败：" + err.Error())
		}
	}

	// 自动加上前缀斜杠
	if params.PatternType == teaconfigs.LocationPatternTypePrefix ||
		params.PatternType == teaconfigs.LocationPatternTypeExact {
		params.Pattern = "/" + strings.TrimLeft(params.Pattern, "/")
	}

	// 匹配条件
	conds, breakCond, err := proxyutils.ParseRequestConds(this.Request, "request")
	if err != nil {
		this.Fail("匹配条件\"" + breakCond.Param + " " + breakCond.Operator + " " + breakCond.Value + "\"校验失败：" + err.Error())
	}

	location.Cond = conds

	// 禁止条件
	denyConds, breakCond, err := proxyutils.ParseRequestConds(this.Request, "deny")
	if err != nil {
		this.Fail("禁止访问条件\"" + breakCond.Param + " " + breakCond.Operator + " " + breakCond.Value + "\"校验失败：" + err.Error())
	}
	location.DenyCond = denyConds
	location.DenyAll = params.DenyAll

	location.SetPattern(params.Pattern, params.PatternType, params.IsCaseInsensitive, params.IsReverse)
	location.On = params.On
	location.IsBreak = params.IsBreak
	location.Name = params.Name
	location.Root = params.Root
	location.URLPrefix = params.URLPrefix
	location.Charset = params.Charset
	location.MaxBodySize = strconv.FormatFloat(params.MaxBodySize, 'f', -1, 64) + params.MaxBodyUnit
	if params.AccessLogIsInherited {
		location.AccessLog = []*teaconfigs.AccessLogConfig{}
	} else {
		location.AccessLog = proxyutils.ParseAccessLogForm(this.Request)
	}
	location.DisableStat = !params.EnableStat
	location.RedirectToHttps = params.RedirectToHttps

	// gzip
	// 这里gzipLevel包括0，因为要指定不压缩
	if params.GzipLevel >= 0 && params.GzipLevel <= 9 {
		minLength := strconv.FormatFloat(params.GzipMinLength, 'f', -1, 64) + params.GzipMinUnit
		gzip := &teaconfigs.GzipConfig{
			Level:     params.GzipLevel,
			MinLength: minLength,
			MimeTypes: params.GzipMimeTypeValues,
		}
		location.Gzip = gzip
	} else {
		location.Gzip = nil
	}

	// 特殊页面
	location.Pages = []*teaconfigs.PageConfig{}
	for index, status := range params.PageStatusList {
		page := teaconfigs.NewPageConfig()
		page.Status = []string{status}
		if index < len(params.PageURLList) {
			page.URL = params.PageURLList[index]
		}
		if index < len(params.PageNewStatusList) {
			page.NewStatus = types.Int(params.PageNewStatusList[index])
			if page.NewStatus < 0 {
				page.NewStatus = 0
			}
		}
		location.AddPage(page)
	}

	if location.Shutdown != nil {
		location.Shutdown.On = params.ShutdownPageOn
		location.Shutdown.URL = params.ShutdownPageURL
		location.Shutdown.Status = params.ShutdownPageStatus
	} else if params.ShutdownPageOn {
		location.Shutdown = teaconfigs.NewShutdownConfig()
		location.Shutdown.On = params.ShutdownPageOn
		location.Shutdown.URL = params.ShutdownPageURL
		location.Shutdown.Status = params.ShutdownPageStatus
	}
	if location.Shutdown != nil && location.Shutdown.On && len(location.Shutdown.URL) == 0 {
		this.FailField("shutdownPageURL", "请输入临时关闭页面文件路径")
	}

	// 首页
	index := []string{}
	for _, i := range params.Index {
		if len(i) > 0 && !lists.ContainsString(index, i) {
			index = append(index, i)
		}
	}
	location.Index = index

	err = server.Save()
	if err != nil {
		this.Fail("保存失败：" + err.Error())
	}

	proxyutils.NotifyChange()
	this.Success()
}
