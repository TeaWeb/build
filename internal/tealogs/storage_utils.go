package tealogs

import (
	"github.com/iwind/TeaGo/maps"
)

// 所有存储引擎列表
func AllStorages() []maps.Map {
	return []maps.Map{
		{
			"name":        "文件",
			"type":        StorageTypeFile,
			"description": "将日志存储在磁盘文件中",
		},
		{
			"name":        "ElasticSearch",
			"type":        StorageTypeES,
			"description": "将日志存储在ElasticSearch中",
		},
		{
			"name":        "MySQL",
			"type":        StorageTypeMySQL,
			"description": "将日志存储在MySQL中",
		},
		{
			"name":        "TCP Socket",
			"type":        StorageTypeTCP,
			"description": "将日志通过TCP套接字输出",
		},
		{
			"name":        "Syslog",
			"type":        StorageTypeSyslog,
			"description": "将日志通过syslog输出，仅支持Linux",
		},
		{
			"name":        "命令行输入流",
			"type":        StorageTypeCommand,
			"description": "启动一个命令通过读取stdin接收日志信息",
		},
	}
}

// 根据类型查找名称
func FindStorageTypeName(storageType string) string {
	for _, m := range AllStorages() {
		if m.GetString("type") == storageType {
			return m.GetString("name")
		}
	}
	return ""
}
