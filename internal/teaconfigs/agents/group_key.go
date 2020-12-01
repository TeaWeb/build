package agents

import (
	"github.com/iwind/TeaGo/rands"
	timeutil "github.com/iwind/TeaGo/utils/time"
)

// Agent组密钥管理
type GroupKey struct {
	Id          string `yaml:"id" json:"id"`
	On          bool   `yaml:"on" json:"on"`
	Name        string `yaml:"name" json:"name"`               // 说明
	Key         string `yaml:"key" json:"key"`                 // 密钥
	DayFrom     string `yaml:"dayFrom" json:"dayFrom"`         // 开始生效日期
	DayTo       string `yaml:"dayTo" json:"dayTo"`             // 结束生效日期
	MaxAgents   int    `yaml:"maxAgents" json:"maxAgents"`     // 可以容纳的Agents最大数量
	CountAgents int    `yaml:"countAgents" json:"countAgents"` // 已使用的Agent
}

// 创建新Key
func NewGroupKey() *GroupKey {
	return &GroupKey{
		Id: rands.HexString(16),
		On: true,
	}
}

// 当前分组Key日期是否可用
func (this *GroupKey) IsDateAvailable() bool {
	today := timeutil.Format("Y-m-d")
	if len(this.DayFrom) > 0 && this.DayFrom > today {
		return false
	}
	if len(this.DayTo) > 0 && this.DayTo < today {
		return false
	}
	return true
}

// 当前分组Key是否可用
func (this *GroupKey) IsAvailable() bool {
	if !this.On {
		return false
	}
	if this.MaxAgents > 0 && this.CountAgents >= this.MaxAgents {
		return false
	}
	return this.IsDateAvailable()
}
