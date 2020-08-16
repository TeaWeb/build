package teastats

import (
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// PV统计
type PVPagePeriodFilter struct {
	CounterFilter
}

func (this *PVPagePeriodFilter) Name() string {
	return "URL PV统计"
}

func (this *PVPagePeriodFilter) Description() string {
	return "单个URL PV统计"
}

func (this *PVPagePeriodFilter) Codes() []string {
	return []string{
		"pv.page.second",
		"pv.page.minute",
		"pv.page.hour",
		"pv.page.day",
		"pv.page.week",
		"pv.page.month",
		"pv.page.year",
	}
}

// 参数说明
func (this *PVPagePeriodFilter) ParamVariables() []*Variable {
	return []*Variable{
		NewVariable("page", "请求URL"),
	}
}

// 统计数据说明
func (this *PVPagePeriodFilter) ValueVariables() []*Variable {
	return []*Variable{
		NewVariable("count", "PV数"),
	}
}

func (this *PVPagePeriodFilter) Indexes() []string {
	return []string{"page"}
}

func (this *PVPagePeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

func (this *PVPagePeriodFilter) Filter(accessLog *accesslogs.AccessLog) {
	contentType := accessLog.SentContentType()
	if !strings.HasPrefix(contentType, "text/html") {
		return
	}
	this.ApplyFilter(accessLog, map[string]string{
		"page": accessLog.RequestPath,
	}, maps.Map{
		"count": 1,
	})
}

func (this *PVPagePeriodFilter) Stop() {
	this.StopFilter()
}
