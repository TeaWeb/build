package tealogs

import "github.com/iwind/TeaGo/maps"

// 存储日志的格式
type StorageFormat = string

const (
	StorageFormatJSON     StorageFormat = "json"
	StorageFormatTemplate StorageFormat = "template"
)

// 所有存储的格式
func AllStorageFormats() []maps.Map {
	return []maps.Map{
		{
			"name":        "JSON",
			"code":        StorageFormatJSON,
			"description": "完整的JSON格式",
		},
		{
			"name":        "模板",
			"code":        StorageFormatTemplate,
			"description": "可以通过使用变量组织一个字符串模板",
		},
	}
}

// 根据代号查找名称
func FindStorageFormatName(code string) string {
	for _, m := range AllStorageFormats() {
		if m.GetString("code") == code {
			return m.GetString("name")
		}
	}
	return ""
}
