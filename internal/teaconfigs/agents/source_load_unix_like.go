// +build !windows

package agents

import "github.com/shirou/gopsutil/load"

// 执行
func (this *LoadSource) Execute(params map[string]string) (value interface{}, err error) {
	stat, err := load.Avg()
	if err != nil || stat == nil {
		return
	}
	value = map[string]interface{}{
		"load1":  stat.Load1,
		"load5":  stat.Load5,
		"load15": stat.Load15,
	}
	return
}
