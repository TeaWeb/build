package shared

import (
	"github.com/iwind/TeaGo/utils/string"
	"regexp"
)

var regexpNamedVariable = regexp.MustCompile("\\${[\\w.-]+}")

// 头部信息定义
type HeaderConfig struct {
	On     bool   `yaml:"on" json:"on"`         // 是否开启
	Id     string `yaml:"id" json:"id"`         // ID
	Name   string `yaml:"name" json:"name"`     // Name
	Value  string `yaml:"value" json:"value"`   // Value
	Always bool   `yaml:"always" json:"always"` // 是否忽略状态码
	Status []int  `yaml:"status" json:"status"` // 支持的状态码

	statusMap    map[int]bool
	hasVariables bool
}

// 获取新Header对象
func NewHeaderConfig() *HeaderConfig {
	return &HeaderConfig{
		On: true,
		Id: stringutil.Rand(16),
	}
}

// 校验
func (this *HeaderConfig) Validate() error {
	this.statusMap = map[int]bool{}
	this.hasVariables = regexpNamedVariable.MatchString(this.Value)

	if this.Status == nil {
		this.Status = []int{}
	}

	for _, status := range this.Status {
		this.statusMap[status] = true
	}

	return nil
}

// 判断是否匹配状态码
func (this *HeaderConfig) Match(statusCode int) bool {
	if !this.On {
		return false
	}

	if this.Always {
		return true
	}

	if this.statusMap != nil {
		_, found := this.statusMap[statusCode]
		return found
	}

	return false
}

// 是否有变量
func (this *HeaderConfig) HasVariables() bool {
	return this.hasVariables
}
