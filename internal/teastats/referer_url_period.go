package teastats

import (
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// 来源URL统计
type RefererURLPeriodFilter struct {
	CounterFilter
}

func (this *RefererURLPeriodFilter) Name() string {
	return "来源URL统计"
}

func (this *RefererURLPeriodFilter) Description() string {
	return "所有请求的来源URL统计"
}

func (this *RefererURLPeriodFilter) Codes() []string {
	return []string{
		"referer.url.second",
		"referer.url.minute",
		"referer.url.hour",
		"referer.url.day",
		"referer.url.week",
		"referer.url.month",
		"referer.url.year",
	}
}

// 参数说明
func (this *RefererURLPeriodFilter) ParamVariables() []*Variable {
	return []*Variable{
		NewVariable("url", "来源URL"),
	}
}

// 统计数据说明
func (this *RefererURLPeriodFilter) ValueVariables() []*Variable {
	return []*Variable{
		NewVariable("count", "请求数"),
	}
}

func (this *RefererURLPeriodFilter) Indexes() []string {
	return []string{"url"}
}

func (this *RefererURLPeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

func (this *RefererURLPeriodFilter) Filter(accessLog *accesslogs.AccessLog) {
	referer := accessLog.GetHeader("Referer")
	if len(referer) == 0 {
		return
	}

	contentType := accessLog.SentContentType()
	if !strings.HasPrefix(contentType, "text/html") {
		return
	}

	this.ApplyFilter(accessLog, map[string]string{
		"url": referer,
	}, maps.Map{
		"count": 1,
	})
}

func (this *RefererURLPeriodFilter) Stop() {
	this.StopFilter()
}
