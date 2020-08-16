package agents

import (
	"github.com/TeaWeb/build/internal/teaconfigs/forms"
	"github.com/TeaWeb/build/internal/teaconfigs/widgets"
)

// 数据源基础定义
type Source struct {
	DataFormat SourceDataFormat `yaml:"dataFormat" json:"dataFormat"` // 数据格式
}

// 代号
func (this *Source) Code() string {
	return ""
}

// 获得数据格式
func (this *Source) DataFormatCode() SourceDataFormat {
	return this.DataFormat
}

// 描述
func (this *Source) Description() string {
	return ""
}

// 校验
func (this *Source) Validate() error {
	return nil
}

// 显示信息
func (this *Source) Presentation() *forms.Presentation {
	return nil
}

// 数据变量定义
func (this *Source) Variables() []*SourceVariable {
	return nil
}

// 阈值
func (this *Source) Thresholds() []*Threshold {
	return nil
}

// 图表
func (this *Source) Charts() []*widgets.Chart {
	return nil
}

// 支持的平台
func (this *Source) Platforms() []string {
	return []string{"darwin", "linux", "windows"}
}

// 停止
func (this *Source) Stop() error {
	return nil
}
