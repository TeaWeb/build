package teastats

import (
	"fmt"
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// 状态码统计
type StatusAllPeriodFilter struct {
	CounterFilter
}

func (this *StatusAllPeriodFilter) Name() string {
	return "状态码统计"
}

func (this *StatusAllPeriodFilter) Description() string {
	return "所有请求的响应状态码统计"
}

// 参数说明
func (this *StatusAllPeriodFilter) ParamVariables() []*Variable {
	return []*Variable{
		NewVariable("status", "HTTP状态码"),
	}
}

// 统计数据说明
func (this *StatusAllPeriodFilter) ValueVariables() []*Variable {
	return []*Variable{
		NewVariable("count", "请求数量"),
	}
}

func (this *StatusAllPeriodFilter) Codes() []string {
	return []string{
		"status.all.second",
		"status.all.minute",
		"status.all.hour",
		"status.all.day",
		"status.all.week",
		"status.all.month",
		"status.all.year",
	}
}

func (this *StatusAllPeriodFilter) Indexes() []string {
	return []string{"status"}
}

func (this *StatusAllPeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

func (this *StatusAllPeriodFilter) Filter(accessLog *accesslogs.AccessLog) {
	this.ApplyFilter(accessLog, map[string]string{
		"status": fmt.Sprintf("%d", accessLog.Status),
	}, maps.Map{
		"count": 1,
	})
}

func (this *StatusAllPeriodFilter) Stop() {
	this.StopFilter()
}
