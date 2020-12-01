package notices

import (
	"errors"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/rands"
	"github.com/iwind/TeaGo/utils/time"
	"time"
)

// 媒介配置定义
type NoticeMediaConfig struct {
	Id          string                 `yaml:"id" json:"id"`
	On          bool                   `yaml:"on" json:"on"`
	Name        string                 `yaml:"name" json:"name"`
	Type        NoticeMediaType        `yaml:"type" json:"type"`
	Options     map[string]interface{} `yaml:"options" json:"options"`
	TimeFrom    string                 `yaml:"timeFrom" json:"timeFrom"`       // 发送的开始时间
	TimeTo      string                 `yaml:"timeTo" json:"timeTo"`           // 发送的结束时间
	RateMinutes int                    `yaml:"rateMinutes" json:"rateMinutes"` // 速率限制之时间范围
	RateCount   int                    `yaml:"rateCount" json:"rateCount"`     // 速率限制之数量
}

// 获取新对象
func NewNoticeMediaConfig() *NoticeMediaConfig {
	return &NoticeMediaConfig{
		On: true,
		Id: rands.HexString(16),
	}
}

// 校验
func (this *NoticeMediaConfig) Validate() error {
	return nil
}

// 取得原始的媒介
func (this *NoticeMediaConfig) Raw() (NoticeMediaInterface, error) {
	m := FindNoticeMediaType(this.Type)
	if m == nil {
		return nil, errors.New("media type '" + this.Type + "' not found")
	}
	instance := m["instance"]
	err := teautils.MapToObjectJSON(this.Options, instance)
	if err != nil {
		return nil, err
	}
	return instance.(NoticeMediaInterface), nil
}

// 是否应该推送
func (this *NoticeMediaConfig) ShouldNotify(countSent int) bool {
	if !this.On {
		return false
	}

	// 时间范围检查
	nowTime := timeutil.Format("H:i:s", time.Now())
	if len(this.TimeFrom) > 0 && nowTime < this.TimeFrom {
		return false
	}
	if len(this.TimeTo) > 0 && this.TimeTo != "00:00:00" && nowTime > this.TimeTo {
		return false
	}

	// 发送频率检查
	if this.RateMinutes > 0 || this.RateCount > 0 {
		if this.RateMinutes <= 0 || this.RateCount <= 0 {
			return false
		}

		// 最近发送的通知频率
		if countSent >= this.RateCount {
			return false
		}
	}

	return true
}
