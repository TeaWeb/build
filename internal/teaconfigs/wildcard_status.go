package teaconfigs

import (
	"fmt"
	"github.com/iwind/TeaGo/types"
	"regexp"
	"strings"
)

// 可能含有x字母的状态码
type WildcardStatus struct {
	StatusInt    int
	StatusRegexp *regexp.Regexp
}

// 获取新对象
func NewWildcardStatus(status string) *WildcardStatus {
	status = regexp.MustCompile("[^0-9x]").ReplaceAllString(status, "")
	if strings.Contains(status, "x") {
		return &WildcardStatus{
			StatusRegexp: regexp.MustCompile("^" + strings.Replace(status, "x", "\\d", -1) + "$"),
		}
	} else {
		return &WildcardStatus{
			StatusInt: types.Int(status),
		}
	}
}

// 判断匹配
func (this *WildcardStatus) Match(status int) bool {
	if this.StatusRegexp != nil {
		return this.StatusRegexp.MatchString(fmt.Sprintf("%d", status))
	}
	return this.StatusInt == status
}
