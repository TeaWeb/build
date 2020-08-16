package teastats

import (
	"fmt"
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// 状态码统计
type StatusPagePeriodFilter struct {
	CounterFilter
}

func (this *StatusPagePeriodFilter) Name() string {
	return "URL状态码统计"
}

func (this *StatusPagePeriodFilter) Description() string {
	return "单个URL响应状态码统计"
}

// 参数说明
func (this *StatusPagePeriodFilter) ParamVariables() []*Variable {
	return []*Variable{
		NewVariable("status", "HTTP状态码"),
		NewVariable("page", "请求URL"),
	}
}

// 统计数据说明
func (this *StatusPagePeriodFilter) ValueVariables() []*Variable {
	return []*Variable{
		NewVariable("count", "请求数量"),
	}
}

// 提供的代码
func (this *StatusPagePeriodFilter) Codes() []string {
	return []string{
		"status.page.second",
		"status.page.minute",
		"status.page.hour",
		"status.page.day",
		"status.page.week",
		"status.page.month",
		"status.page.year",
	}
}

func (this *StatusPagePeriodFilter) Indexes() []string {
	return []string{"status", "page"}
}

// 启动
func (this *StatusPagePeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

// 筛选
func (this *StatusPagePeriodFilter) Filter(accessLog *accesslogs.AccessLog) {
	this.ApplyFilter(accessLog, map[string]string{
		"status": fmt.Sprintf("%d", accessLog.Status),
		"page":   accessLog.RequestPath,
	}, maps.Map{
		"count": 1,
	})
}

// 停止
func (this *StatusPagePeriodFilter) Stop() {
	this.StopFilter()
}
