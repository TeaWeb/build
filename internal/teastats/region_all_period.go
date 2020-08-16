package teastats

import (
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// 区域统计
type RegionAllPeriodFilter struct {
	CounterFilter
}

func (this *RegionAllPeriodFilter) Name() string {
	return "区域统计"
}

func (this *RegionAllPeriodFilter) Description() string {
	return "所有请求的区域统计"
}

func (this *RegionAllPeriodFilter) Codes() []string {
	return []string{
		"region.all.second",
		"region.all.minute",
		"region.all.hour",
		"region.all.day",
		"region.all.week",
		"region.all.month",
		"region.all.year",
	}
}

// 参数说明
func (this *RegionAllPeriodFilter) ParamVariables() []*Variable {
	return []*Variable{
		NewVariable("region", "国家或地区"),
	}
}

// 统计数据说明
func (this *RegionAllPeriodFilter) ValueVariables() []*Variable {
	return []*Variable{
		NewVariable("countReq", "请求数"),
		NewVariable("countPV", "PV数量"),
		NewVariable("countUV", "UV数量"),
		NewVariable("countIP", "IP数量"),
	}
}

func (this *RegionAllPeriodFilter) Indexes() []string {
	return []string{"region"}
}

func (this *RegionAllPeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

func (this *RegionAllPeriodFilter) Filter(accessLog *accesslogs.AccessLog) {
	if accessLog.Extend == nil {
		return
	}
	if len(accessLog.Extend.Geo.Region) == 0 {
		return
	}

	// 中国特区
	if accessLog.Extend.Geo.Region == "台湾" {
		accessLog.Extend.Geo.Region = "中国台湾"
	} else if accessLog.Extend.Geo.Region == "香港" {
		accessLog.Extend.Geo.Region = "中国香港"
	} else if accessLog.Extend.Geo.Region == "澳门" {
		accessLog.Extend.Geo.Region = "中国澳门"
	}

	countPV := 0
	countUV := 0
	countIP := 0

	if strings.HasPrefix(accessLog.SentContentType(), "text/html") {
		countPV++
	}

	if this.CheckNewUV(accessLog, accessLog.Extend.Geo.Region) {
		countUV = 1
	}

	if this.CheckNewIP(accessLog, accessLog.Extend.Geo.Region) {
		countIP = 1
	}

	this.ApplyFilter(accessLog, map[string]string{
		"region": accessLog.Extend.Geo.Region,
	}, maps.Map{
		"countReq": 1,
		"countPV":  countPV,
		"countUV":  countUV,
		"countIP":  countIP,
	})
}

func (this *RegionAllPeriodFilter) Stop() {
	this.StopFilter()
}
