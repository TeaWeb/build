package teastats

import (
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// 登陆页统计
type LandingPagePeriodFilter struct {
	CounterFilter
}

func (this *LandingPagePeriodFilter) Name() string {
	return "登陆页统计"
}

func (this *LandingPagePeriodFilter) Description() string {
	return "登陆页统计"
}

func (this *LandingPagePeriodFilter) Codes() []string {
	return []string{
		"landing.page.second",
		"landing.page.minute",
		"landing.page.hour",
		"landing.page.day",
		"landing.page.week",
		"landing.page.month",
		"landing.page.year",
	}
}

// 参数说明
func (this *LandingPagePeriodFilter) ParamVariables() []*Variable {
	return []*Variable{
		NewVariable("page", "请求URL"),
	}
}

// 统计数据说明
func (this *LandingPagePeriodFilter) ValueVariables() []*Variable {
	return []*Variable{
		NewVariable("count", "请求数"),
	}
}

func (this *LandingPagePeriodFilter) Indexes() []string {
	return []string{"page"}
}

func (this *LandingPagePeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

func (this *LandingPagePeriodFilter) Filter(accessLog *accesslogs.AccessLog) {
	contentType := accessLog.SentContentType()
	if !strings.HasPrefix(contentType, "text/html") {
		return
	}
	uid, ok := accessLog.Cookie["TeaUID"]
	if ok && len(uid) > 0 {
		return
	}
	this.ApplyFilter(accessLog, map[string]string{
		"page": accessLog.RequestPath,
	}, maps.Map{
		"count": 1,
	})
}

func (this *LandingPagePeriodFilter) Stop() {
	this.StopFilter()
}
