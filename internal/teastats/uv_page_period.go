package teastats

import (
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// UV统计
type UVPagePeriodFilter struct {
	CounterFilter
}

func (this *UVPagePeriodFilter) Name() string {
	return "URL UV统计"
}

func (this *UVPagePeriodFilter) Description() string {
	return "单个URL UV统计"
}

func (this *UVPagePeriodFilter) Codes() []string {
	return []string{
		"uv.page.second",
		"uv.page.minute",
		"uv.page.hour",
		"uv.page.day",
		"uv.page.week",
		"uv.page.month",
		"uv.page.year",
	}
}

// 参数说明
func (this *UVPagePeriodFilter) ParamVariables() []*Variable {
	return []*Variable{
		NewVariable("page", "请求URL"),
	}
}

// 统计数据说明
func (this *UVPagePeriodFilter) ValueVariables() []*Variable {
	return []*Variable{
		NewVariable("count", "UV数量"),
	}
}

func (this *UVPagePeriodFilter) Indexes() []string {
	return []string{"page"}
}

func (this *UVPagePeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

func (this *UVPagePeriodFilter) Filter(accessLog *accesslogs.AccessLog) {
	if !this.CheckNewUV(accessLog, accessLog.RequestPath) {
		return
	}

	this.ApplyFilter(accessLog, map[string]string{
		"page": accessLog.RequestPath,
	}, maps.Map{
		"count": 1,
	})
}

func (this *UVPagePeriodFilter) Stop() {
	this.StopFilter()
}
