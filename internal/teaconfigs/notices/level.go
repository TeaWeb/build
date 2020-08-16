package notices

import (
	"github.com/iwind/TeaGo/maps"
)

// 通知级别类型
type NoticeLevel = uint8

// 通知级别常量
const (
	NoticeLevelNone    = NoticeLevel(0)
	NoticeLevelInfo    = NoticeLevel(1)
	NoticeLevelWarning = NoticeLevel(2)
	NoticeLevelError   = NoticeLevel(3)
	NoticeLevelSuccess = NoticeLevel(4)
)

// 所有的通知级别
func AllNoticeLevels() []maps.Map {
	return []maps.Map{
		{
			"name":        "信息",
			"description": "通常是一般的指标数据信息",
			"code":        NoticeLevelInfo,
			"bgcolor":     "#f8ffff",
			"color":       "#276f86",
		},
		{
			"name":        "警告",
			"description": "可能会发生异常的警告信息",
			"code":        NoticeLevelWarning,
			"bgcolor":     "#fffaf3",
			"color":       "#573a08",
		},
		{
			"name":        "错误",
			"description": "发生了错误信息",
			"code":        NoticeLevelError,
			"bgcolor":     "#fff6f6",
			"color":       "#9f3a38",
		},
		{
			"name":        "成功",
			"description": "某个任务处理成功之后的通知信息",
			"code":        NoticeLevelSuccess,
			"bgcolor":     "#fcfff5",
			"color":       "#2c662d",
		},
	}
}

// 获取通知级别名称
func FindNoticeLevelName(level NoticeLevel) string {
	for _, l := range AllNoticeLevels() {
		if l["code"] == level {
			return l["name"].(string)
		}
	}
	return "信息"
}

// 获取通知级别信息
func FindNoticeLevel(level NoticeLevel) maps.Map {
	for _, l := range AllNoticeLevels() {
		if l["code"] == level {
			return l
		}
	}
	return FindNoticeLevel(NoticeLevelInfo)
}

// 判断是否为失败级别
func IsFailureLevel(level NoticeLevel) bool {
	return level == NoticeLevelWarning || level == NoticeLevelError
}
