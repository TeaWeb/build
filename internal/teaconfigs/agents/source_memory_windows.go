// +build windows

package agents

import (
	"errors"
	"github.com/shirou/gopsutil/mem"
	"runtime"
	"syscall"
	"unsafe"
)

// 内存状态
// https://msdn.microsoft.com/en-us/library/windows/desktop/aa366589(v=vs.85).aspx
type memStatusEx struct {
	dwLength         uint32
	dwMemoryLoad     uint32
	ullTotalPhys     uint64
	ullAvailPhys     uint64
	ullTotalPageFile int64
	ullAvailPageFile int64
	ullTotalVirtual  uint64
	ullAvailVirtual  uint64
}

// 执行
func (this *MemorySource) Execute(params map[string]string) (value interface{}, err error) {
	stat, err := mem.VirtualMemory()
	if err != nil {
		return
	}

	swap, err := mem.SwapMemory()
	if err != nil {
		return
	} else if swap.Total == 0 {
		// 修复老系统（比如windows 2003/windows xp）上的虚拟内存计算为0的问题
		result, err := this.sysTotalMemory()
		if err == nil && result != nil && result.ullTotalVirtual >= result.ullAvailVirtual {
			swap.Total = result.ullTotalVirtual
			swap.Used = result.ullTotalVirtual - result.ullAvailVirtual
			swap.Free = result.ullAvailVirtual
			swap.UsedPercent = float64(swap.Used*100) / float64(swap.Total)
		}
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

func (this *MemorySource) sysTotalMemory() (*memStatusEx, error) {
	kernel32, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		return nil, err
	}

	globalMemoryStatusEx, err := kernel32.FindProc("GlobalMemoryStatusEx")
	if err != nil {
		return nil, err
	}
	msx := &memStatusEx{
		dwLength: 64,
	}
	r, _, _ := globalMemoryStatusEx.Call(uintptr(unsafe.Pointer(msx)))
	if r == 0 {
		return nil, errors.New("retrieve data failed")
	}
	return msx, nil
}
