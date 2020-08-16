package forms

import (
	"html"
	"net/http"
	"strings"
)

// 元素接口
type ElementInterface interface {
	Compose() string
	Super() *Element
	ApplyRequest(req *http.Request) (value interface{}, skip bool, err error)
}

// 元素公共定义
type Element struct {
	ClassType  string            `yaml:"classType" json:"classType"`   // 字段类型
	Attrs      map[string]string `yaml:"attrs" json:"attrs"`           // 附加属性
	Namespace  string            `yaml:"namespace" json:"namespace"`   // 命名空间
	Code       string            `yaml:"code" json:"code"`             // 字段值代号
	Title      string            `yaml:"title" json:"title"`           // 标题
	Subtitle   string            `yaml:"subtitle" json:"subtitle"`     // 副标题
	IsRequired bool              `yaml:"isRequired" json:"isRequired"` // 是否为必填，并没有实际的约束作用，只是用来在字段左边标记星号
	IsComposed bool              `yaml:"isComposed" json:"isComposed"` // 是否已经组合，组合后的不需要再次组合

	Comment      string      `yaml:"comment" json:"comment"`           // 注释
	ValidateCode string      `yaml:"validateCode" json:"validateCode"` // 值校验代码
	InitCode     string      `yaml:"initCode" json:"initCode"`         // 值初始化代码
	Value        interface{} `yaml:"value" json:"value"`               // 字段值
	Javascript   string      `yaml:"javascript" json:"javascript"`     // 附加的Javascript代码
	CSS          string      `yaml:"css" json:"css"`                   // 附加的CSS代码
}

func (this *Element) ComposeAttrs(attrs map[string]string) string {
	composedAttrs := map[string]string{}
	for k, v := range this.Attrs {
		composedAttrs[k] = v
	}
	for k, v := range attrs {
		composedAttrs[k] = v
	}

	list := []string{}
	for name, value := range composedAttrs {
		list = append(list, name+"=\""+html.EscapeString(value)+"\"")
	}
	return strings.Join(list, " ")
}

func (this *Element) Attr(name string, value string) {
	if this.Attrs == nil {
		this.Attrs = map[string]string{}
	}
	this.Attrs[name] = value
}
