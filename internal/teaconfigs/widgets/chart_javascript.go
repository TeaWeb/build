package widgets

import (
	"github.com/iwind/TeaGo/utils/string"
	"regexp"
	"strings"
)

var chartInstanceRegexp = regexp.MustCompile("(\\w+)\\s*=\\s*new\\s+charts\\s*\\.\\s*\\w+\\s*\\(\\s*\\)")

// Javascript
type JavascriptChart struct {
	Code string `yaml:"code" json:"code"`
}

func (this *JavascriptChart) AsJavascript(options map[string]interface{}) (code string, err error) {
	code = this.Code

	code = chartInstanceRegexp.ReplaceAllStringFunc(code, func(s string) string {
		varName := s[:strings.Index(s, "=")]
		return s + "; " + varName + ".options = " + stringutil.JSONEncode(options)
	})

	return code, nil
}
