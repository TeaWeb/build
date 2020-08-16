package notices

import (
	"github.com/iwind/TeaGo/maps"
	"github.com/iwind/TeaGo/types"
	"net/http"
	"regexp"
	"strings"
)

var paramNamedVariable = regexp.MustCompile("\\${[$\\w. \t-]+}")

// 通知条目配置，用于某个项目的细项通知配置
type Item struct {
	On        bool              `yaml:"on" json:"on"` // 是否开启
	Level     NoticeLevel       `yaml:"level" json:"level"`
	Receivers []*NoticeReceiver `yaml:"receivers" json:"receivers"`
	Subject   string            `yaml:"subject" json:"subject"`
	Body      string            `yaml:"body" json:"body"`
}

// 获取新对象
func NewItem(level NoticeLevel) *Item {
	return &Item{
		Level: level,
	}
}

// 从请求中获取新对象
func NewItemFromRequest(req *http.Request, name string) *Item {
	item := &Item{}

	{
		v := req.Form.Get(name + "NoticeOn")
		if v == "1" {
			item.On = true
		}
	}

	{
		v := req.Form.Get(name + "NoticeLevel")
		item.Level = types.Uint8(v)
	}

	item.Subject = req.Form.Get(name + "NoticeSubject")
	item.Body = req.Form.Get(name + "NoticeBody")

	return item
}

// 取得替换变量后的标题
func (this *Item) FormatSubject(vars maps.Map) string {
	return paramNamedVariable.ReplaceAllStringFunc(this.Subject, func(s string) string {
		varName := s[2 : len(s)-1]
		varName = strings.Replace(varName, " ", "", -1)
		varName = strings.Replace(varName, "\t", "", -1)
		if vars != nil {
			v, ok := vars[varName]
			if ok {
				return types.String(v)
			}
		}
		return s
	})
}

// 取得替换变量后的内容
func (this *Item) FormatBody(vars maps.Map) string {
	return paramNamedVariable.ReplaceAllStringFunc(this.Body, func(s string) string {
		varName := s[2 : len(s)-1]
		varName = strings.Replace(varName, " ", "", -1)
		varName = strings.Replace(varName, "\t", "", -1)
		if vars != nil {
			v, ok := vars[varName]
			if ok {
				return types.String(v)
			}
		}
		return s
	})
}
