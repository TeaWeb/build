// +build windows

package agents

import (
	"context"
	"github.com/shirou/gopsutil/cpu"
	"math"
	"sync"
	"time"
)

type WindowsLoadValue struct {
	Timestamp int64
	Value     int
}

var windowsLoadValues = []*WindowsLoadValue{}
var windowsLoadLocker sync.Mutex

// 执行
func (this *LoadSource) Execute(params map[string]string) (value interface{}, err error) {
	timestamp := time.Now().Unix()

	currentLoad := 0
	info, err := cpu.ProcInfo()
	if err == nil && len(info) > 0 && info[0].ProcessorQueueLength < 1000 {
		currentLoad = int(info[0].ProcessorQueueLength)
	}

	// 删除15分钟之前的数据
	windowsLoadLocker.Lock()
	result := []*WindowsLoadValue{}
	for _, v := range windowsLoadValues {
		if timestamp-v.Timestamp > 15*60 {
			continue
		}
		result = append(result, v)
	}
	result = append(result, &WindowsLoadValue{
		Timestamp: timestamp,
		Value:     currentLoad,
	})
	windowsLoadValues = result

	total1 := 0
	count1 := 0
	total5 := 0
	count5 := 0
	total15 := 0
	count15 := 0
	for _, v := range result {
		if timestamp-v.Timestamp <= 60 {
			total1 += v.Value
			count1 ++
		}

		if timestamp-v.Timestamp <= 300 {
			total5 += v.Value
			count5 ++
		}

		total15 += v.Value
		count15 ++
	}

	load1 := float64(0)
	load5 := float64(0)
	load15 := float64(0)
	if count1 > 0 {
		load1 = math.Round(float64(total1*100)/float64(count1)) / 100
	}
	if count5 > 0 {
		load5 = math.Round(float64(total5*100)/float64(count5)) / 100
	}
	if count15 > 0 {
		load15 = math.Round(float64(total15*100)/float64(count15)) / 100
	}

	windowsLoadLocker.Unlock()

	// 在老Windows上不显示错误
	if err == context.DeadlineExceeded {
		err = nil
	}

	return map[string]interface{}{
		"load1":  load1,
		"load5":  load5,
		"load15": load15,
	}, err
}
