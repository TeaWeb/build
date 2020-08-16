package teautils

import (
	"regexp"
	"strings"
)

var whitespaceReg = regexp.MustCompile(`\s+`)

// 关键词匹配
func MatchKeyword(source, keyword string) bool {
	if len(keyword) == 0 {
		return false
	}

	pieces := whitespaceReg.Split(keyword, -1)
	source = strings.ToLower(source)
	for _, piece := range pieces {
		if strings.Index(source, strings.ToLower(piece)) > -1 {
			return true
		}
	}

	return false
}
