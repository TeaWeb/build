package teastats

import (
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/TeaWeb/build/internal/teawaf"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// WAF拦截统计
type WAFBlockAllPeriodFilter struct {
	CounterFilter
}

func (this *WAFBlockAllPeriodFilter) Name() string {
	return "WAF拦截类型统计"
}

func (this *WAFBlockAllPeriodFilter) Description() string {
	return "所有WAF拦截类型统计统计"
}

// 参数说明
func (this *WAFBlockAllPeriodFilter) ParamVariables() []*Variable {
	return []*Variable{
		NewVariable("wafId", "WAF ID"),
		NewVariable("ruleSetId", "规则集ID"),
		NewVariable("ruleSetName", "规则集名称"),
	}
}

// 统计数据说明
func (this *WAFBlockAllPeriodFilter) ValueVariables() []*Variable {
	return []*Variable{
		NewVariable("count", "匹配的数量"),
	}
}

func (this *WAFBlockAllPeriodFilter) Codes() []string {
	return []string{
		"waf.block.all.second",
		"waf.block.all.minute",
		"waf.block.all.hour",
		"waf.block.all.day",
		"waf.block.all.week",
		"waf.block.all.month",
		"waf.block.all.year",
	}
}

func (this *WAFBlockAllPeriodFilter) Indexes() []string {
	return []string{"wafId", "ruleSetId", "ruleSetName"}
}

func (this *WAFBlockAllPeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

func (this *WAFBlockAllPeriodFilter) Filter(accessLog *accesslogs.AccessLog) {
	if accessLog.Attrs == nil {
		return
	}

	wafAction, ok := accessLog.Attrs["waf_action"]
	if !ok {
		return
	}
	if wafAction != teawaf.ActionBlock {
		return
	}

	wafId, ok := accessLog.Attrs["waf_id"]
	if !ok {
		return
	}

	ruleSetId, ok := accessLog.Attrs["waf_ruleset"]
	if !ok {
		return
	}

	ruleSetName, ok := accessLog.Attrs["waf_ruleset_name"]
	if !ok {
		return
	}

	this.ApplyFilter(accessLog, map[string]string{
		"wafId":       wafId,
		"ruleSetId":   ruleSetId,
		"ruleSetName": ruleSetName,
	}, maps.Map{
		"count": 1,
	})
}

func (this *WAFBlockAllPeriodFilter) Stop() {
	this.StopFilter()
}
