package api

import "time"

// api数据量限制
type APIDataLimit struct {
	Max      uint   `yaml:"max" json:"max"`           // 最大数据量 TODO
	Total    uint   `yaml:"total" json:"total"`       // 数据量 TODO
	Duration string `yaml:"duration" json:"duration"` // 数据限制间隔 TODO

	duration time.Duration
}

// 校验
func (this *APIDataLimit) Validate() error {
	return nil
}
