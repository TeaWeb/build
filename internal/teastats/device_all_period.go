package teastats

import (
	"github.com/TeaWeb/build/internal/tealogs/accesslogs"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
	"strings"
)

// 设备统计
type DeviceAllPeriodFilter struct {
	CounterFilter
}

func (this *DeviceAllPeriodFilter) Name() string {
	return "设备统计"
}

func (this *DeviceAllPeriodFilter) Description() string {
	return "所有请求的设备统计"
}

func (this *DeviceAllPeriodFilter) Codes() []string {
	return []string{
		"device.all.second",
		"device.all.minute",
		"device.all.hour",
		"device.all.day",
		"device.all.week",
		"device.all.month",
		"device.all.year",
	}
}

// 参数说明
func (this *DeviceAllPeriodFilter) ParamVariables() []*Variable {
	return []*Variable{
		NewVariable("family", "设备名"),
		NewVariable("model", "型号"),
	}
}

// 统计数据说明
func (this *DeviceAllPeriodFilter) ValueVariables() []*Variable {
	return []*Variable{
		NewVariable("countReq", "请求数"),
		NewVariable("countPV", "PV数量"),
		NewVariable("countUV", "UV数量"),
		NewVariable("countIP", "IP数量"),
	}
}

func (this *DeviceAllPeriodFilter) Indexes() []string {
	return []string{"family", "model"}
}

func (this *DeviceAllPeriodFilter) Start(queue *Queue, code string) {
	if queue == nil {
		logs.Println("stat queue should be specified for '" + code + "'")
		return
	}
	this.queue = queue
	this.queue.Index(this.Indexes())
	this.StartFilter(code, code[strings.LastIndex(code, ".")+1:])
}

func (this *DeviceAllPeriodFilter) Filter(accessLog *accesslogs.AccessLog) {
	if accessLog.Extend == nil {
		return
	}
	if len(accessLog.Extend.Client.Device.Family) == 0 {
		return
	}

	countPV := 0
	countUV := 0
	countIP := 0

	if strings.HasPrefix(accessLog.SentContentType(), "text/html") {
		countPV++
	}

	if this.CheckNewUV(accessLog, accessLog.Extend.Client.Device.Family+"_"+accessLog.Extend.Client.Device.Model) {
		countUV = 1
	}

	if this.CheckNewIP(accessLog, accessLog.Extend.Client.Device.Family+"_"+accessLog.Extend.Client.Device.Model) {
		countIP = 1
	}

	this.ApplyFilter(accessLog, map[string]string{
		"family": accessLog.Extend.Client.Device.Family,
		"model":  accessLog.Extend.Client.Device.Model,
	}, maps.Map{
		"countReq": 1,
		"countPV":  countPV,
		"countUV":  countUV,
		"countIP":  countIP,
	})
}

func (this *DeviceAllPeriodFilter) Stop() {
	this.StopFilter()
}
