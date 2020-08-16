package teastats

import (
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// 流量统计
type TrafficPagePeriodFilter struct {
	CounterFilter
}

func (this *TrafficPagePeriodFilter) Name() string {
	return "URL流量统计"
}

func (this *TrafficPagePeriodFilter) Description() string {
	return "单个URL流量统计"
}

// 参数说明
func (this *TrafficPagePeriodFilter) ParamVariables() []*Variable {
	return []*Variable{
		NewVariable("page", "请求URL"),
	}
}

// 统计数据说明
func (this *TrafficPagePeriodFilter) ValueVariables() []*Variable {
	return []*Variable{
		NewVariable("bytes", "流量（字节）"),
	}
}

func (this *TrafficPagePeriodFilter) Codes() []string {
	return []string{
		"traffic.page.second",
		"traffic.page.minute",
		"traffic.page.hour",
		"traffic.page.day",
		"traffic.page.week",
		"traffic.page.month",
		"traffic.page.year",
	}
}

func (this *TrafficPagePeriodFilter) Indexes() []string {
	return []string{"page"}
}

func (this *TrafficPagePeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

func (this *TrafficPagePeriodFilter) Filter(accessLog *accesslogs.AccessLog) {
	this.ApplyFilter(accessLog, map[string]string{
		"page": accessLog.RequestPath,
	}, maps.Map{
		"bytes": accessLog.BytesSent,
	})
}

func (this *TrafficPagePeriodFilter) Stop() {
	this.StopFilter()
}
