package teastats

import (
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// PV统计
type PVAllPeriodFilter struct {
	CounterFilter
}

func (this *PVAllPeriodFilter) Name() string {
	return "PV统计"
}

func (this *PVAllPeriodFilter) Description() string {
	return "所有请求的PV统计"
}

func (this *PVAllPeriodFilter) Codes() []string {
	return []string{
		"pv.all.second",
		"pv.all.minute",
		"pv.all.hour",
		"pv.all.day",
		"pv.all.week",
		"pv.all.month",
		"pv.all.year",
	}
}

// 参数说明
func (this *PVAllPeriodFilter) ParamVariables() []*Variable {
	return []*Variable{}
}

// 统计数据说明
func (this *PVAllPeriodFilter) ValueVariables() []*Variable {
	return []*Variable{
		NewVariable("count", "PV数"),
	}
}

func (this *PVAllPeriodFilter) Indexes() []string {
	return []string{}
}

func (this *PVAllPeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

func (this *PVAllPeriodFilter) Filter(accessLog *accesslogs.AccessLog) {
	contentType := accessLog.SentContentType()
	if !strings.HasPrefix(contentType, "text/html") {
		return
	}
	this.ApplyFilter(accessLog, nil, maps.Map{
		"count": 1,
	})
}

func (this *PVAllPeriodFilter) Stop() {
	this.StopFilter()
}
