package teastats

import (
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// IP统计
type IPPagePeriodFilter struct {
	CounterFilter
}

func (this *IPPagePeriodFilter) Name() string {
	return "URL IP统计"
}

func (this *IPPagePeriodFilter) Description() string {
	return "单个URL IP统计"
}

func (this *IPPagePeriodFilter) Codes() []string {
	return []string{
		"ip.page.second",
		"ip.page.minute",
		"ip.page.hour",
		"ip.page.day",
		"ip.page.week",
		"ip.page.month",
		"ip.page.year",
	}
}

// 参数说明
func (this *IPPagePeriodFilter) ParamVariables() []*Variable {
	return []*Variable{
		NewVariable("page", "请求URL"),
	}
}

// 统计数据说明
func (this *IPPagePeriodFilter) ValueVariables() []*Variable {
	return []*Variable{
		NewVariable("count", "请求数"),
	}
}

func (this *IPPagePeriodFilter) Indexes() []string {
	return []string{"page"}
}

func (this *IPPagePeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

func (this *IPPagePeriodFilter) Filter(accessLog *accesslogs.AccessLog) {
	if !this.CheckNewIP(accessLog, accessLog.RequestPath) {
		return
	}

	this.ApplyFilter(accessLog, map[string]string{
		"page": accessLog.RequestPath,
	}, maps.Map{
		"count": 1,
	})
}

func (this *IPPagePeriodFilter) Stop() {
	this.StopFilter()
}
