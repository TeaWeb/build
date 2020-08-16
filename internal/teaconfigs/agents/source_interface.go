package agents

import (
	"github.com/TeaWeb/build/internal/teaconfigs/forms"
	"github.com/TeaWeb/build/internal/teaconfigs/widgets"
)

// 数据源接口定义
type SourceInterface interface {
	// 名称
	Name() string

	// 代号
	Code() string

	// 描述
	Description() string

	// 校验
	Validate() error

	// 执行
	Execute(params map[string]string) (value interface{}, err error)

	// 获得数据格式
	DataFormatCode() SourceDataFormat

	// 表单信息
	Form() *forms.Form

	// 显示信息
	Presentation() *forms.Presentation

	// 数据变量定义
	Variables() []*SourceVariable

	// 阈值
	Thresholds() []*Threshold

	// 图表
	Charts() []*widgets.Chart

	// 支持的平台
	Platforms() []string

	// 停止
	Stop() error
}
