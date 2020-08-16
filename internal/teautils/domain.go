package teautils

import (
	"github.com/iwind/TeaGo/logs"
	"github.com/iwind/TeaGo/utils/string"
	"strings"
)

// 从一组规则中匹配域名
// 支持的格式：example.com, www.example.com, .example.com, *.example.com, ~(\d+).example.com
// 更多参考：http://nginx.org/en/docs/http/ngx_http_core_module.html#server_name
func MatchDomains(patterns []string, domain string) (isMatched bool) {
	if len(patterns) == 0 {
		return
	}
	for _, pattern := range patterns {
		if matchDomain(pattern, domain) {
			return true
		}
	}
	return
}

// 匹配单个域名规则
func matchDomain(pattern string, domain string) (isMatched bool) {
	if len(pattern) == 0 {
		return
	}

	// 正则表达式
	if pattern[0] == '~' {
		reg, err := stringutil.RegexpCompile(strings.TrimSpace(pattern[1:]))
		if err != nil {
			logs.Error(err)
			return false
		}
		return reg.MatchString(domain)
	}

	if pattern[0] == '.' {
		return strings.HasSuffix(domain, pattern)
	}

	// 其他匹配
	patternPieces := strings.Split(pattern, ".")
	domainPieces := strings.Split(domain, ".")
	if len(patternPieces) != len(domainPieces) {
		return
	}
	isMatched = true
	for index, patternPiece := range patternPieces {
		if patternPiece == "" || patternPiece == "*" || patternPiece == domainPieces[index] {
			continue
		}
		isMatched = false
		break
	}
	return isMatched
}
