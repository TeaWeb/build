package shared

import (
	"sync"
	"time"
)

// API控制策略
type AccessPolicy struct {
	locker    sync.Mutex
	isChanged bool

	// 流量控制
	Traffic struct {
		On    bool `yaml:"on" json:"on"` // 是否开启
		Total struct {
			On    bool  `yaml:"on" json:"on"`       // 是否开启
			Total int64 `yaml:"total" json:"total"` // 总量
			Used  int64 `yaml:"used" json:"used"`   // 已使用量
		} `yaml:"total" json:"total"` // 总量控制
		Second struct {
			On       bool  `yaml:"on" json:"on"`             // 是否开启
			Total    int64 `yaml:"total" json:"total"`       // 总量
			Duration int64 `yaml:"duration" json:"duration"` // 时间长度
			FromTime int64 `yaml:"fromTime" json:"fromTime"` // 开始时间，也是结束时间
			Used     int64 `yaml:"used" json:"used"`         // 已使用量
		} `yaml:"second" json:"second"`
		Minute struct {
			On       bool  `yaml:"on" json:"on"`             // 是否开启
			Total    int64 `yaml:"total" json:"total"`       // 总量
			Duration int64 `yaml:"duration" json:"duration"` // 时间长度
			FromTime int64 `yaml:"fromTime" json:"fromTime"` // 开始时间
			ToTime   int64 `yaml:"toTime" json:"toTime"`     // 结束时间
			Used     int64 `yaml:"used" json:"used"`         // 已使用量
		} `yaml:"minute" json:"minute"`
		Hour struct {
			On       bool  `yaml:"on" json:"on"`             // 是否开启
			Total    int64 `yaml:"total" json:"total"`       // 总量
			Duration int64 `yaml:"duration" json:"duration"` // 时间长度
			FromTime int64 `yaml:"fromTime" json:"fromTime"` // 开始时间
			ToTime   int64 `yaml:"toTime" json:"toTime"`     // 结束时间
			Used     int64 `yaml:"used" json:"used"`         // 已使用量
		} `yaml:"hour" json:"hour"`
		Day struct {
			On       bool  `yaml:"on" json:"on"`             // 是否开启
			Total    int64 `yaml:"total" json:"total"`       // 总量
			Duration int64 `yaml:"duration" json:"duration"` // 时间长度
			FromTime int64 `yaml:"fromTime" json:"fromTime"` // 开始时间
			ToTime   int64 `yaml:"toTime" json:"toTime"`     // 结束时间
			Used     int64 `yaml:"used" json:"used"`         // 已使用量
		} `yaml:"day" json:"day"`
		Month struct {
			On       bool  `yaml:"on" json:"on"`             // 是否开启
			Total    int64 `yaml:"total" json:"total"`       // 总量
			Duration int64 `yaml:"duration" json:"duration"` // 时间长度
			FromTime int64 `yaml:"fromTime" json:"fromTime"` // 开始时间
			ToTime   int64 `yaml:"toTime" json:"toTime"`     // 结束时间
			Used     int64 `yaml:"used" json:"used"`         // 已使用量
		} `yaml:"month" json:"month"`
	} `yaml:"traffic" json:"traffic"` // 流量控制

	// 访问控制
	Access AccessConfig `yaml:"access" json:"access"` // 访问控制
}

// 获取新对象
func NewAccessPolicy() *AccessPolicy {
	return &AccessPolicy{}
}

// 校验
func (this *AccessPolicy) Validate() error {
	err := this.Access.Validate()
	if err != nil {
		return err
	}
	return nil
}

// 检查权限
func (this *AccessPolicy) AllowAccess(ip string) bool {
	// Access
	if this.Access.On {
		// deny
		if this.Access.DenyOn {
			for _, client := range this.Access.Deny {
				if !client.On {
					continue
				}
				if client.Match(ip) {
					return false
				}
			}
		}

		// allow
		if this.Access.AllowOn {
			for _, client := range this.Access.Allow {
				if !client.On {
					continue
				}
				if client.Match(ip) {
					return true
				}
			}
			return false
		}
	}
	return true
}

// 检查流量
func (this *AccessPolicy) AllowTraffic() (reason string, allowed bool) {
	if !this.Traffic.On {
		return "", true
	}

	this.locker.Lock()
	defer this.locker.Unlock()

	now := time.Now()
	timestamp := now.Unix()

	// total
	if this.Traffic.Total.On {
		if this.Traffic.Total.Used >= this.Traffic.Total.Total {
			return "total", false
		}
	}

	// second
	if this.Traffic.Second.On {
		if this.Traffic.Second.Duration <= 0 || this.Traffic.Second.Total <= 0 {
			return "second", false
		}
		if timestamp-this.Traffic.Second.FromTime < this.Traffic.Second.Duration && this.Traffic.Second.Used >= this.Traffic.Second.Total {
			return "second", false
		}
	}

	// minute
	if this.Traffic.Minute.On {
		if this.Traffic.Minute.Duration <= 0 || this.Traffic.Minute.Total <= 0 {
			return "minute", false
		}
		if timestamp >= this.Traffic.Minute.FromTime && timestamp < this.Traffic.Minute.ToTime && this.Traffic.Minute.Used >= this.Traffic.Minute.Total {
			return "minute", false
		}
	}

	// hour
	if this.Traffic.Hour.On && this.Traffic.Hour.Duration > 0 {
		if this.Traffic.Hour.Duration <= 0 || this.Traffic.Hour.Total <= 0 {
			return "hour", false
		}
		if timestamp >= this.Traffic.Hour.FromTime && timestamp < this.Traffic.Hour.ToTime && this.Traffic.Hour.Used >= this.Traffic.Hour.Total {
			return "hour", false
		}
	}

	// day
	if this.Traffic.Day.On && this.Traffic.Day.Duration > 0 {
		if this.Traffic.Day.Duration <= 0 || this.Traffic.Day.Total <= 0 {
			return "day", false
		}
		if timestamp >= this.Traffic.Day.FromTime && timestamp < this.Traffic.Day.ToTime && this.Traffic.Day.Used >= this.Traffic.Day.Total {
			return "day", false
		}
	}

	// month
	if this.Traffic.Month.On && this.Traffic.Month.Duration > 0 {
		if this.Traffic.Month.Duration <= 0 || this.Traffic.Month.Total <= 0 {
			return "month", false
		}
		if timestamp >= this.Traffic.Month.FromTime && timestamp < this.Traffic.Month.ToTime && this.Traffic.Month.Used >= this.Traffic.Month.Total {
			return "month", false
		}
	}

	this.IncreaseTraffic()

	return "", true
}

// 增加流量
func (this *AccessPolicy) IncreaseTraffic() {
	if !this.Traffic.On {
		return
	}

	this.isChanged = true

	now := time.Now()
	timestamp := now.Unix()
	year := now.Year()
	month := now.Month()
	day := now.Day()
	hour := now.Hour()
	minute := now.Minute()

	// total
	if this.Traffic.Total.On {
		this.Traffic.Total.Used++
	}

	// second
	if this.Traffic.Second.On && this.Traffic.Second.Duration > 0 {
		if timestamp-this.Traffic.Second.FromTime < this.Traffic.Second.Duration {
			this.Traffic.Second.Used++
		} else {
			this.Traffic.Second.FromTime = timestamp
			this.Traffic.Second.Used = 1
		}
	}

	// minute
	if this.Traffic.Minute.On && this.Traffic.Minute.Duration > 0 {
		if timestamp >= this.Traffic.Minute.FromTime && timestamp < this.Traffic.Minute.ToTime {
			this.Traffic.Minute.Used++
		} else {
			this.Traffic.Minute.Used = 1
			fromTime := time.Date(year, month, day, hour, minute, 0, 0, time.Local)
			this.Traffic.Minute.FromTime = fromTime.Unix()
			this.Traffic.Minute.ToTime = fromTime.Add(time.Duration(this.Traffic.Minute.Duration) * time.Minute).Unix()
		}
	}

	// hour
	if this.Traffic.Hour.On && this.Traffic.Hour.Duration > 0 {
		if timestamp >= this.Traffic.Hour.FromTime && timestamp < this.Traffic.Hour.ToTime {
			this.Traffic.Hour.Used++
		} else {
			this.Traffic.Hour.Used = 1
			fromTime := time.Date(year, month, day, hour, 0, 0, 0, time.Local)
			this.Traffic.Hour.FromTime = fromTime.Unix()
			this.Traffic.Hour.ToTime = fromTime.Add(time.Duration(this.Traffic.Hour.Duration) * time.Hour).Unix()
		}
	}

	// day
	if this.Traffic.Day.On && this.Traffic.Day.Duration > 0 {
		if timestamp >= this.Traffic.Day.FromTime && timestamp < this.Traffic.Day.ToTime {
			this.Traffic.Day.Used++
		} else {
			this.Traffic.Day.Used = 1
			fromTime := time.Date(year, month, day, 0, 0, 0, 0, time.Local)
			this.Traffic.Day.FromTime = fromTime.Unix()
			this.Traffic.Day.ToTime = fromTime.AddDate(0, 0, int(this.Traffic.Day.Duration)).Unix()
		}
	}

	// month
	if this.Traffic.Month.On && this.Traffic.Month.Duration > 0 {
		if timestamp >= this.Traffic.Month.FromTime && timestamp < this.Traffic.Month.ToTime {
			this.Traffic.Month.Used++
		} else {
			this.Traffic.Month.Used = 1
			fromTime := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
			this.Traffic.Month.FromTime = fromTime.Unix()
			this.Traffic.Month.ToTime = fromTime.AddDate(0, int(this.Traffic.Month.Duration), 0).Unix()
		}
	}
}

// 判断是否改变
func (this *AccessPolicy) IsChanged() bool {
	return this.isChanged
}

// 设置已完成改变
func (this *AccessPolicy) FinishChange() {
	this.isChanged = false
}
