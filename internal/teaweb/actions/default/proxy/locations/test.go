package locations

import (
	"github.com/TeaWeb/build/internal/teaconfigs"
	"github.com/TeaWeb/build/internal/teaproxy"
	"github.com/TeaWeb/build/internal/teaweb/actions/default/proxy/proxyutils"
	"github.com/iwind/TeaGo/actions"
	"net/http"
)

type TestAction actions.Action

// 测试
func (this *TestAction) Run(params struct {
	Pattern           string
	PatternType       int
	IsReverse         bool
	IsCaseInsensitive bool
	TestingPath       string
}) {
	location := teaconfigs.NewLocation()
	location.Pattern = params.Pattern
	location.SetPattern(params.Pattern, params.PatternType, params.IsCaseInsensitive, params.IsReverse)

	// 匹配条件
	conds, breakCond, err := proxyutils.ParseRequestConds(this.Request, "request")
	if err != nil {
		this.Fail("匹配条件\"" + breakCond.Param + " " + breakCond.Operator + " " + breakCond.Value + "\"校验失败：" + err.Error())
	}
	location.Cond = conds

	err = location.Validate()
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

	mapping, ok := location.Match(rawReq.URL.Path, func(source string) string {
		if req == nil {
			return source
		} else {
			return req.Format(source)
		}
	})
	if ok {
		this.Data["mapping"] = mapping
		this.Success()
	} else {
		this.Data["mapping"] = nil
		this.Fail()
	}
}
