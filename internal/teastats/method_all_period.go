package teastats

import (
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// 请求方法统计
type MethodAllPeriodFilter struct {
	CounterFilter
}

func (this *MethodAllPeriodFilter) Name() string {
	return "请求方法统计"
}

func (this *MethodAllPeriodFilter) Description() string {
	return "所有请求的请求方法统计"
}

func (this *MethodAllPeriodFilter) Codes() []string {
	return []string{
		"method.all.second",
		"method.all.minute",
		"method.all.hour",
		"method.all.day",
		"method.all.week",
		"method.all.month",
		"method.all.year",
	}
}

// 参数说明
func (this *MethodAllPeriodFilter) ParamVariables() []*Variable {
	return []*Variable{
		NewVariable("method", "请求方法（大写）"),
	}
}

// 统计数据说明
func (this *MethodAllPeriodFilter) ValueVariables() []*Variable {
	return []*Variable{
		NewVariable("count", "请求数"),
	}
}

func (this *MethodAllPeriodFilter) Indexes() []string {
	return []string{"method"}
}

func (this *MethodAllPeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

func (this *MethodAllPeriodFilter) Filter(accessLog *accesslogs.AccessLog) {
	this.ApplyFilter(accessLog, map[string]string{
		"method": accessLog.RequestMethod,
	}, maps.Map{
		"count": 1,
	})
}

func (this *MethodAllPeriodFilter) Stop() {
	this.StopFilter()
}
