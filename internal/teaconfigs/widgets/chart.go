package widgets

import (
	"errors"
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/rands"
)

// Chart接口
type ChartInterface interface {
	AsJavascript(options map[string]interface{}) (code string, err error)
}

// Chart定义
type Chart struct {
	Id                string                 `yaml:"id" json:"id"`
	On                bool                   `yaml:"on" json:"on"`
	Name              string                 `yaml:"name" json:"name"`
	Description       string                 `yaml:"description" json:"description"`
	Columns           uint8                  `yaml:"columns" json:"columns"`                     // 列
	Type              string                 `yaml:"type" json:"type"`                           // 类型
	Options           map[string]interface{} `yaml:"options" json:"options"`                     // 选项
	Requirements      []string               `yaml:"requirements" json:"requirements"`           // 绘制chart需要的特征
	SupportsTimeRange bool                   `yaml:"supportsTimeRange" json:"supportsTimeRange"` // 是否支持时间范围查询
}

// 获取新对象
func NewChart() *Chart {
	return &Chart{
		On: true,
		Id: rands.HexString(16),
	}
}

// 校验
func (this *Chart) Validate() error {
	return nil
}

// 转换为具体对象
func (this *Chart) AsObject() (ChartInterface, error) {
	for _, chart := range AllChartTypes {
		if chart["code"] != this.Type {
			continue
		}
		instance, ok := chart["instance"].(ChartInterface)
		if ok {
			err := teautils.MapToObjectJSON(this.Options, instance)
			return instance, err
		} else {
			return nil, errors.New("chart instance should implement ChartInterface: '" + this.Type + "'")
		}
	}

	return nil, errors.New("invalid chart type '" + this.Type + "'")
}
