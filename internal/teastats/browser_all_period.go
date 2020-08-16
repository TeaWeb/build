package teastats

import (
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// 浏览器统计
type BrowserAllPeriodFilter struct {
	CounterFilter
}

func (this *BrowserAllPeriodFilter) Name() string {
	return "浏览器统计"
}

func (this *BrowserAllPeriodFilter) Description() string {
	return "所有请求的浏览器统计"
}

func (this *BrowserAllPeriodFilter) Codes() []string {
	return []string{
		"browser.all.second",
		"browser.all.minute",
		"browser.all.hour",
		"browser.all.day",
		"browser.all.week",
		"browser.all.month",
		"browser.all.year",
	}
}

// 参数说明
func (this *BrowserAllPeriodFilter) ParamVariables() []*Variable {
	return []*Variable{
		NewVariable("family", "浏览器名称"),
		NewVariable("major", "浏览器主版本"),
	}
}

// 统计数据说明
func (this *BrowserAllPeriodFilter) ValueVariables() []*Variable {
	return []*Variable{
		NewVariable("countReq", "请求数"),
		NewVariable("countPV", "PV数量"),
		NewVariable("countUV", "UV数量"),
		NewVariable("countIP", "IP数量"),
	}
}

func (this *BrowserAllPeriodFilter) Indexes() []string {
	return []string{"family", "major"}
}

func (this *BrowserAllPeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

func (this *BrowserAllPeriodFilter) Filter(accessLog *accesslogs.AccessLog) {
	if accessLog.Extend == nil {
		return
	}
	if len(accessLog.Extend.Client.Browser.Family) == 0 {
		return
	}

	countPV := 0
	countUV := 0
	countIP := 0

	if strings.HasPrefix(accessLog.SentContentType(), "text/html") {
		countPV++
	}

	if this.CheckNewUV(accessLog, accessLog.Extend.Client.Browser.Family+"_"+accessLog.Extend.Client.Browser.Major) {
		countUV = 1
	}

	if this.CheckNewIP(accessLog, accessLog.Extend.Client.Browser.Family+"_"+accessLog.Extend.Client.Browser.Major) {
		countIP = 1
	}

	this.ApplyFilter(accessLog, map[string]string{
		"family": accessLog.Extend.Client.Browser.Family,
		"major":  accessLog.Extend.Client.Browser.Major,
	}, maps.Map{
		"countReq": 1,
		"countPV":  countPV,
		"countUV":  countUV,
		"countIP":  countIP,
	})
}

func (this *BrowserAllPeriodFilter) Stop() {
	this.StopFilter()
}
