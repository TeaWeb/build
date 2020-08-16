// +build !windows

package agents

import (
	"github.com/shirou/gopsutil/mem"
	"runtime"
)

// 执行
func (this *MemorySource) Execute(params map[string]string) (value interface{}, err error) {
	stat, err := mem.VirtualMemory()
	if err != nil {
		return
	}

	swap, err := mem.SwapMemory()
	if err != nil {
		return
	}

	// 重新计算内存
	if stat.Total > 0 {
		if runtime.GOOS == "darwin" || runtime.GOOS == "linux" || runtime.GOOS == "freebsd" {
			stat.Used = stat.Total - stat.Free - stat.Buffers - stat.Cached
			stat.UsedPercent = float64(stat.Used) * 100 / float64(stat.Total)
		}
	}

	value = map[string]interface{}{
		"usage": map[string]interface{}{
			"virtualUsed":    float64(stat.Used) / 1024 / 1024 / 1024,
			"virtualPercent": stat.UsedPercent,
			"virtualTotal":   float64(stat.Total) / 1024 / 1024 / 1024,
			"virtualFree":    float64(stat.Free) / 1024 / 1024 / 1024,
			"virtualWired":   float64(stat.Wired) / 1024 / 1024 / 1024,
			"virtualBuffers": float64(stat.Buffers) / 1024 / 1024 / 1024,
			"virtualCached":  float64(stat.Cached) / 1024 / 1024 / 1024,
			"swapUsed":       float64(swap.Used) / 1024 / 1024 / 1024,
			"swapPercent":    swap.UsedPercent,
			"swapTotal":      float64(swap.Total) / 1024 / 1024 / 1024,
			"swapFree":       float64(swap.Free) / 1024 / 1024 / 1024,
		},
	}

	return
}
