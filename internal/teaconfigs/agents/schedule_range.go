package agents

// Schedule时间范围
type ScheduleRangeConfig struct {
	Every bool `yaml:"every" json:"every"`
	From  int  `yaml:"from" json:"from"`
	To    int  `yaml:"to" json:"to"`
	Step  int  `yaml:"step" json:"step"`
	Value int  `yaml:"value" json:"value"`
}

// 获取新对象
func NewScheduleRangeConfig() *ScheduleRangeConfig {
	return &ScheduleRangeConfig{
		Every: false,
		From:  -1,
		To:    -1,
		Step:  -1,
		Value: -1,
	}
}
