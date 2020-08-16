package teaconfigs

import "github.com/iwind/TeaGo/maps"

// 匹配类型
type LocationPatternType = int

// 内置的匹配类型定义
const (
	LocationPatternTypePrefix = 1
	LocationPatternTypeExact  = 2
	LocationPatternTypeRegexp = 3
)

// 取得所有的匹配类型信息
func AllLocationPatternTypes() []maps.Map {
	return []maps.Map{
		{
			"name":        "匹配前缀",
			"type":        LocationPatternTypePrefix,
			"description": "带有此前缀的路径才会被匹配",
		},
		{
			"name":        "精准匹配",
			"type":        LocationPatternTypeExact,
			"description": "带此路径完全一样的路径才会被匹配",
		},
		{
			"name":        "正则表达式匹配",
			"type":        LocationPatternTypeRegexp,
			"description": "通过正则表达式来匹配路径，<a href=\"http://teaos.cn/doc/regexp/Regexp.md\" target=\"_blank\">正则表达式语法 &raquo;</a>",
		},
	}
}

// 查找单个匹配类型信息
func FindLocationPatternType(patternType int) maps.Map {
	for _, t := range AllLocationPatternTypes() {
		if t["type"] == patternType {
			return t
		}
	}
	return nil
}

// 查找单个匹配类型名称
func FindLocationPatternTypeName(patternType int) string {
	t := FindLocationPatternType(patternType)
	if t == nil {
		return ""
	}
	return t["name"].(string)
}
