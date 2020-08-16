package teacache

import "github.com/iwind/TeaGo/maps"

// 所有的缓存配置
func AllCacheTypes() []maps.Map {
	return []maps.Map{
		{
			"name":        "内存",
			"type":        "memory",
			"description": "将缓存数据存储在内存中",
		},
		{
			"name":        "文件",
			"type":        "file",
			"description": "将缓存数据存储在本地文件中",
		},
		{
			"name":        "Redis",
			"type":        "redis",
			"description": "将缓存数据存储在Redis服务中",
		},
		{
			"name":        "LevelDB",
			"type":        "leveldb",
			"description": "将缓存数据存储在LevelDB本地数据库中",
		},
	}
}

// 查找类型名称
func FindTypeName(typeCode string) string {
	for _, m := range AllCacheTypes() {
		if m.GetString("type") == typeCode {
			return m.GetString("name")
		}
	}
	return ""
}

// 查找类型信息
func FindType(typeCode string) maps.Map {
	for _, m := range AllCacheTypes() {
		if m.GetString("type") == typeCode {
			return m
		}
	}
	return maps.Map{}
}
