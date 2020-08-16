package teaconfigs

import "github.com/iwind/TeaGo/maps"

// 调度算法配置
type SchedulingConfig struct {
	Code    string   `yaml:"code" json:"code"`       // 类型
	Options maps.Map `yaml:"options" json:"options"` // 选项
}

// 获取新对象
func NewSchedulingConfig() *SchedulingConfig {
	return &SchedulingConfig{}
}
