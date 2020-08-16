package agents

import "github.com/TeaWeb/build/internal/teaconfigs"

// 看板图表定义
type BoardChart struct {
	AppId   string `yaml:"appId" json:"appId"`
	ItemId  string `yaml:"itemId" json:"itemId"`
	ChartId string `yaml:"chartId" json:"chartId"`

	Name     string              `yaml:"name" json:"name"`
	TimeType string              `yaml:"timeType" json:"timeType"` // default, past, range
	TimePast teaconfigs.TimePast `yaml:"timePast" json:"timePast"` // 时间范围
	DayFrom  string              `yaml:"dayFrom" json:"dayFrom"`
	DayTo    string              `yaml:"dayTo" json:"dayTo"`
}
