package forms

import (
	stringutil "github.com/iwind/TeaGo/utils/string"
	"net/http"
)

// 单值列表
type SingleValueList struct {
	Element `yaml:",inline"`

	ValueName string `yaml:"valueName" json:"valueName"`
}

// 获取新对象
func NewSingleValueList(title string, subtitle string) *SingleValueList {
	return &SingleValueList{
		Element: Element{
			Title:    title,
			Subtitle: subtitle,
		},
	}
}

// 获取父级类型
func (this *SingleValueList) Super() *Element {
	return &this.Element
}

// 组合
func (this *SingleValueList) Compose() string {
	model := this.Namespace + "_" + this.Code + "_single_values"
	this.Javascript = `
this.` + model + ` = ` + stringutil.JSONEncode(this.Value) + `;`

	attrs := this.ComposeAttrs(map[string]string{
		"value-name": this.ValueName,
		"prefix":     this.Code,
		":values":    model,
	})
	return `<single-value-list ` + attrs + `></single-value-list>`
}

// 获取值
func (this *SingleValueList) ApplyRequest(req *http.Request) (value interface{}, skip bool, err error) {
	values, ok := req.Form[this.Code+"Values"]
	if !ok {
		value = []string{}
		return
	}

	value = values
	return
}
