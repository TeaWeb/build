package agents

import (
	"github.com/TeaWeb/build/internal/teautils"
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/maps"
)

// 动作
type ActionInterface interface {
	// 名称
	Name() string

	// 代号
	Code() string

	// 描述
	Description() string

	// 校验
	Validate() error

	// 执行
	Run(params map[string]string) (result string, err error)

	// 获取简要信息
	Summary() maps.Map
}

// 获取所有的数据源信息
func AllActions() []maps.Map {
	result := []maps.Map{}
	for _, i := range []ActionInterface{NewScriptAction()} {
		summary := i.Summary()
		summary["instance"] = i
		result = append(result, summary)
	}
	return result
}

// 查找单个数据源信息
func FindAction(code string) maps.Map {
	for _, summary := range AllActions() {
		if summary["code"] == code {
			return summary
		}
	}
	return nil
}

// 查找单个数据源实例
func FindActionInstance(code string, options map[string]interface{}) ActionInterface {
	for _, summary := range AllActions() {
		if summary["code"] == code {
			instance := summary["instance"].(ActionInterface)
			if options != nil {
				err := teautils.MapToObjectJSON(options, instance)
				if err != nil {
					logs.Error(err)
				}
			}
			return instance
		}
	}
	return nil
}
