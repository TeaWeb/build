package teastats

import (
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// IP统计
type IPAllPeriodFilter struct {
	CounterFilter
}

func (this *IPAllPeriodFilter) Name() string {
	return "IP数统计"
}

func (this *IPAllPeriodFilter) Description() string {
	return "所有请求的IP数统计"
}

func (this *IPAllPeriodFilter) Codes() []string {
	return []string{
		"ip.all.second",
		"ip.all.minute",
		"ip.all.hour",
		"ip.all.day",
		"ip.all.week",
		"ip.all.month",
		"ip.all.year",
	}
}

// 参数说明
func (this *IPAllPeriodFilter) ParamVariables() []*Variable {
	return []*Variable{}
}

// 统计数据说明
func (this *IPAllPeriodFilter) ValueVariables() []*Variable {
	return []*Variable{
		NewVariable("count", "请求数"),
	}
}

func (this *IPAllPeriodFilter) Indexes() []string {
	return []string{}
}

func (this *IPAllPeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

func (this *IPAllPeriodFilter) Filter(accessLog *accesslogs.AccessLog) {
	if !this.CheckNewIP(accessLog, "") {
		return
	}

	this.ApplyFilter(accessLog, nil, maps.Map{
		"count": 1,
	})
}

func (this *IPAllPeriodFilter) Stop() {
	this.StopFilter()
}
