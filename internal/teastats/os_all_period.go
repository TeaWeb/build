package teastats

import (
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// 操作系统统计
type OSAllPeriodFilter struct {
	CounterFilter
}

func (this *OSAllPeriodFilter) Name() string {
	return "操作系统统计"
}

func (this *OSAllPeriodFilter) Description() string {
	return "所有请求的操作系统统计"
}

func (this *OSAllPeriodFilter) Codes() []string {
	return []string{
		"os.all.second",
		"os.all.minute",
		"os.all.hour",
		"os.all.day",
		"os.all.week",
		"os.all.month",
		"os.all.year",
	}
}

// 参数说明
func (this *OSAllPeriodFilter) ParamVariables() []*Variable {
	return []*Variable{
		NewVariable("family", "操作系统名称"),
		NewVariable("major", "操作系统主版本"),
	}
}

// 统计数据说明
func (this *OSAllPeriodFilter) ValueVariables() []*Variable {
	return []*Variable{
		NewVariable("countReq", "请求数"),
		NewVariable("countPV", "PV数量"),
		NewVariable("countUV", "UV数量"),
		NewVariable("countIP", "IP数量"),
	}
}

func (this *OSAllPeriodFilter) Indexes() []string {
	return []string{"family", "major"}
}

func (this *OSAllPeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

func (this *OSAllPeriodFilter) Filter(accessLog *accesslogs.AccessLog) {
	if accessLog.Extend == nil {
		return
	}
	if len(accessLog.Extend.Client.OS.Family) == 0 {
		return
	}

	countPV := 0
	countUV := 0
	countIP := 0

	if strings.HasPrefix(accessLog.SentContentType(), "text/html") {
		countPV++
	}

	if this.CheckNewUV(accessLog, accessLog.Extend.Client.OS.Family+"_"+accessLog.Extend.Client.OS.Major) {
		countUV = 1
	}

	if this.CheckNewIP(accessLog, accessLog.Extend.Client.OS.Family+"_"+accessLog.Extend.Client.OS.Major) {
		countIP = 1
	}

	this.ApplyFilter(accessLog, map[string]string{
		"family": accessLog.Extend.Client.OS.Family,
		"major":  accessLog.Extend.Client.OS.Major,
	}, maps.Map{
		"countReq": 1,
		"countPV":  countPV,
		"countUV":  countUV,
		"countIP":  countIP,
	})
}

func (this *OSAllPeriodFilter) Stop() {
	this.StopFilter()
}
