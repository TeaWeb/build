package rewrite

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaproxy"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"net/http"
	"regexp"
)

type TestAction actions.Action

// 匹配测试
func (this *TestAction) Run(params struct {
	Pattern      string
	Replace      string
	ProxyId      string
	TargetType   string
	RedirectMode string
	TestingPath  string
	Must         *actions.Must
}) {
	params.Must.
		Field("pattern", params.Pattern).
		Require("请输入匹配规则").
		Expect(func() (message string, success bool) {
			_, err := regexp.Compile(params.Pattern)
			if err != nil {
				return "匹配规则错误：" + err.Error(), false
			}
			return "", true
		})

	rewriteRule := teaconfigs.NewRewriteRule()
	rewriteRule.On = true
	rewriteRule.Pattern = params.Pattern
	if params.TargetType == "url" {
		rewriteRule.Replace = params.Replace
	} else {
		rewriteRule.Replace = "proxy://" + params.ProxyId + params.Replace
	}
	if len(params.RedirectMode) > 0 {
		rewriteRule.AddFlag(params.RedirectMode, nil)
	}

	// 匹配条件
	conds, breakCond, err := proxyutils.ParseRequestConds(this.Request, "request")
	if err != nil {
		this.Fail("匹配条件\"" + breakCond.Param + " " + breakCond.Operator + " " + breakCond.Value + "\"校验失败：" + err.Error())
	}
	rewriteRule.Cond = conds

	err = rewriteRule.Validate()
	if err != nil {
		this.Fail("校验失败：" + err.Error())
	}

	rawReq, err := http.NewRequest(http.MethodGet, params.TestingPath, nil)
	if err != nil {
		this.Fail("请输入正确的URL")
	}
	req := teaproxy.NewRequest(rawReq)
	req.SetURI(params.TestingPath)
	req.SetHost(rawReq.Host)
	req.SetRawScheme(rawReq.URL.Scheme)
	replace, mapping, ok := rewriteRule.Match(rawReq.URL.Path, func(source string) string {
		if req == nil {
			return source
		} else {
			return req.Format(source)
		}
	})
	if ok {
		this.Data["replace"] = replace
		this.Data["mapping"] = mapping
		this.Success()
	} else {
		this.Fail()
	}
}
