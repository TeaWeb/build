package agents

import (
	"github.com/iwind/TeaGo/lists"
)

// 定时列表
type ScheduleRangeList struct {
	Every  bool
	Ranges []*ScheduleRangeConfig
}

// 获取新对象
func NewScheduleRangeList() *ScheduleRangeList {
	return &ScheduleRangeList{}
}

// 转换
func (this *ScheduleRangeList) Convert(ranges []*ScheduleRangeConfig) {
	if len(ranges) == 0 {
		this.Every = true
		return
	}
	for _, r := range ranges {
		if r.Every {
			this.Every = true
			return
		}

		if r.Value > -1 {
			this.Ranges = append(this.Ranges, r)
		} else if r.From > -1 && r.To > -1 && r.Step > -1 {
			this.Ranges = append(this.Ranges, r)
		}
	}
	lists.Sort(this.Ranges, func(i int, j int) bool {
		from1 := this.Ranges[i].Value
		if from1 < 0 {
			from1 = this.Ranges[i].From
		}

		from2 := this.Ranges[j].Value
		if from2 < 0 {
			from2 = this.Ranges[j].From
		}
		return from1 < from2
	})
}

// 下一个数值
func (this *ScheduleRangeList) Next(current int, defaultValue int) int {
	if this.Every {
		return current
	}
	if len(this.Ranges) > 0 {
		for _, r := range this.Ranges {
			if r.Value > -1 && r.Value >= current {
				return r.Value
			}

			if r.From > -1 && r.To > -1 && r.Step > -1 && r.To >= current {
				if r.Step <= 0 {
					r.Step = 1
				}
				if current <= r.From {
					return r.From
				}
				mode := (current - r.From) % r.Step
				if mode == 0 {
					return current
				}
				next := current + r.Step - mode
				if next <= r.To {
					return next
				}
			}
		}
	}

	return defaultValue
}
