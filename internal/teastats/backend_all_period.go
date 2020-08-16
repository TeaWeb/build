package teastats

import (
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// 后端统计
type BackendAllPeriodFilter struct {
	CounterFilter
}

func (this *BackendAllPeriodFilter) Name() string {
	return "后端请求统计"
}

func (this *BackendAllPeriodFilter) Description() string {
	return "计算某个后端服务器请求次数"
}

func (this *BackendAllPeriodFilter) Codes() []string {
	return []string{
		"backend.all.second",
		"backend.all.minute",
		"backend.all.hour",
		"backend.all.day",
		"backend.all.week",
		"backend.all.month",
		"backend.all.year",
	}
}

// 参数说明
func (this *BackendAllPeriodFilter) ParamVariables() []*Variable {
	return []*Variable{
		NewVariable("backend", "后端服务器ID"),
	}
}

// 统计数据说明
func (this *BackendAllPeriodFilter) ValueVariables() []*Variable {
	return []*Variable{
		NewVariable("count", "请求数"),
	}
}

func (this *BackendAllPeriodFilter) Indexes() []string {
	return []string{"backend"}
}

func (this *BackendAllPeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

func (this *BackendAllPeriodFilter) Filter(accessLog *accesslogs.AccessLog) {
	if len(accessLog.BackendId) == 0 {
		return
	}
	this.ApplyFilter(accessLog, map[string]string{
		"backend": accessLog.BackendId,
	}, maps.Map{
		"count": 1,
	})
}

func (this *BackendAllPeriodFilter) Stop() {
	this.StopFilter()
}
