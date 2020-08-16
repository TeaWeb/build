package teastats

import (
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// 耗时统计
type CostPagePeriodFilter struct {
	CounterFilter
}

func (this *CostPagePeriodFilter) Name() string {
	return "URL耗时统计"
}

func (this *CostPagePeriodFilter) Description() string {
	return "单个URL耗时统计"
}

func (this *CostPagePeriodFilter) Codes() []string {
	return []string{
		"cost.page.second",
		"cost.page.minute",
		"cost.page.hour",
		"cost.page.day",
		"cost.page.week",
		"cost.page.month",
		"cost.page.year",
	}
}

// 参数说明
func (this *CostPagePeriodFilter) ParamVariables() []*Variable {
	return []*Variable{
		NewVariable("page", "请求URL"),
	}
}

// 统计数据说明
func (this *CostPagePeriodFilter) ValueVariables() []*Variable {
	return []*Variable{
		NewVariable("countReq", "请求数"),
		NewVariable("cost", "耗时"),
	}
}

func (this *CostPagePeriodFilter) Indexes() []string {
	return []string{"page"}
}

func (this *CostPagePeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.IncreaseFunc = func(value maps.Map, inc maps.Map) maps.Map {
		if inc == nil {
			return value
		}
		if value == nil {
			countReq := inc.GetInt64("countReq")
			cost := inc.GetFloat64("cost")
			value = maps.Map{
				"countReq": countReq,
				"cost":     cost / float64(countReq),
			}
		} else {
			totalReq := value.GetInt64("countReq")
			totalCost := value.GetFloat64("cost") * float64(totalReq)

			countReq := inc.GetInt64("countReq")
			cost := inc.GetFloat64("cost")

			value = maps.Map{
				"countReq": totalReq + countReq,
				"cost":     (totalCost + cost) / float64(totalReq+countReq),
			}
		}

		return value
	}
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

func (this *CostPagePeriodFilter) Filter(accessLog *accesslogs.AccessLog) {
	this.ApplyFilter(accessLog, map[string]string{
		"page": accessLog.RequestPath,
	}, maps.Map{
		"countReq": 1,
		"cost":     accessLog.RequestTime,
	})
}

func (this *CostPagePeriodFilter) Stop() {
	this.StopFilter()
}
