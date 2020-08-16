package teastats

import (
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// 耗时统计
type CostAllPeriodFilter struct {
	CounterFilter
}

func (this *CostAllPeriodFilter) Name() string {
	return "耗时统计"
}

func (this *CostAllPeriodFilter) Description() string {
	return "所有请求的耗时统计"
}

func (this *CostAllPeriodFilter) Codes() []string {
	return []string{
		"cost.all.second",
		"cost.all.minute",
		"cost.all.hour",
		"cost.all.day",
		"cost.all.week",
		"cost.all.month",
		"cost.all.year",
	}
}

// 参数说明
func (this *CostAllPeriodFilter) ParamVariables() []*Variable {
	return []*Variable{}
}

// 统计数据说明
func (this *CostAllPeriodFilter) ValueVariables() []*Variable {
	return []*Variable{
		NewVariable("countReq", "请求数"),
		NewVariable("cost", "耗时（秒）"),
	}
}

func (this *CostAllPeriodFilter) Indexes() []string {
	return []string{}
}

func (this *CostAllPeriodFilter) Start(queue *Queue, code string) {
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

func (this *CostAllPeriodFilter) Filter(accessLog *accesslogs.AccessLog) {
	this.ApplyFilter(accessLog, nil, maps.Map{
		"countReq": 1,
		"cost":     accessLog.RequestTime,
	})
}

func (this *CostAllPeriodFilter) Stop() {
	this.StopFilter()
}
