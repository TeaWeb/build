package teastats

import (
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// Fastcgi请求统计
type FastcgiAllPeriodFilter struct {
	CounterFilter
}

func (this *FastcgiAllPeriodFilter) Name() string {
	return "Fastcgi请求统计"
}

func (this *FastcgiAllPeriodFilter) Description() string {
	return "Fastcgi请求统计"
}

func (this *FastcgiAllPeriodFilter) Codes() []string {
	return []string{
		"fastcgi.all.second",
		"fastcgi.all.minute",
		"fastcgi.all.hour",
		"fastcgi.all.day",
		"fastcgi.all.week",
		"fastcgi.all.month",
		"fastcgi.all.year",
	}
}

// 参数说明
func (this *FastcgiAllPeriodFilter) ParamVariables() []*Variable {
	return []*Variable{
		NewVariable("fastcgi", "Fastcgi ID"),
	}
}

// 统计数据说明
func (this *FastcgiAllPeriodFilter) ValueVariables() []*Variable {
	return []*Variable{
		NewVariable("count", "请求数"),
	}
}

func (this *FastcgiAllPeriodFilter) Indexes() []string {
	return []string{"fastcgi"}
}

func (this *FastcgiAllPeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

func (this *FastcgiAllPeriodFilter) Filter(accessLog *accesslogs.AccessLog) {
	if len(accessLog.FastcgiId) == 0 {
		return
	}
	this.ApplyFilter(accessLog, map[string]string{
		"fastcgi": accessLog.FastcgiId,
	}, maps.Map{
		"count": 1,
	})
}

func (this *FastcgiAllPeriodFilter) Stop() {
	this.StopFilter()
}
