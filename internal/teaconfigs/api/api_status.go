package api

import "github.com/iwind/TeaGo/maps"

// API状态常量
const (
	APIStatusTypeNormal  = "normal"
	APIStatusTypeSuccess = "success"
	APIStatusTypeWarning = "warning"
	APIStatusTypeFailure = "failure"
	APIStatusTypeError   = "error"
)

// API状态定义
type APIStatus struct {
	Code        string   `yaml:"code" json:"code"`               // 代码
	Description string   `yaml:"description" json:"description"` // 描述
	Groups      []string `yaml:"groups" json:"groups"`           // 分组
	Versions    []string `yaml:"versions" json:"versions"`       // 版本
	Type        string   `yaml:"type" json:"type"`               // 类型
}

// 获取新对象
func NewAPIStatus() *APIStatus {
	return &APIStatus{}
}

// 所有的状态
func AllStatusTypes() []maps.Map {
	return []maps.Map{
		{
			"name": "成功",
			"code": APIStatusTypeSuccess,
		},
		{
			"name": "警告",
			"code": APIStatusTypeWarning,
		},
		{
			"name": "失败",
			"code": APIStatusTypeFailure,
		},
		{
			"name": "错误",
			"code": APIStatusTypeError,
		},
		{
			"name": "默认",
			"code": APIStatusTypeNormal,
		},
	}
}

// 获取当前的状态类型名称
func (this *APIStatus) TypeName() string {
	for _, m := range AllStatusTypes() {
		if m["code"].(string) == this.Type {
			return m["name"].(string)
		}
	}
	return ""
}
