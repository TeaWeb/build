package teastats

import (
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// UV统计
type UVAllPeriodFilter struct {
	CounterFilter
}

func (this *UVAllPeriodFilter) Name() string {
	return "UV统计"
}

func (this *UVAllPeriodFilter) Description() string {
	return "所有请求的UV统计"
}

func (this *UVAllPeriodFilter) Codes() []string {
	return []string{
		"uv.all.second",
		"uv.all.minute",
		"uv.all.hour",
		"uv.all.day",
		"uv.all.week",
		"uv.all.month",
		"uv.all.year",
	}
}

// 参数说明
func (this *UVAllPeriodFilter) ParamVariables() []*Variable {
	return []*Variable{}
}

// 统计数据说明
func (this *UVAllPeriodFilter) ValueVariables() []*Variable {
	return []*Variable{
		NewVariable("count", "UV数"),
	}
}

func (this *UVAllPeriodFilter) Indexes() []string {
	return []string{}
}

func (this *UVAllPeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

func (this *UVAllPeriodFilter) Filter(accessLog *accesslogs.AccessLog) {
	if !this.CheckNewUV(accessLog, "") {
		return
	}

	this.ApplyFilter(accessLog, nil, maps.Map{
		"count": 1,
	})
}

func (this *UVAllPeriodFilter) Stop() {
	this.StopFilter()
}
