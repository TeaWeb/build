package forms

import (
	"github.com/TeaWeb/build/internal/teaconfigs/shared"
	"github.com/iwind/TeaGo/maps"
	stringutil "github.com/iwind/TeaGo/utils/string"
	"html"
	"net/http"
)

// HTTP参数组件
type HTTPBox struct {
	Element `yaml:",inline"`
}

// 获取新对象
func NewHTTPBox(title string, subtitle string) *HTTPBox {
	return &HTTPBox{
		Element: Element{
			Title:      title,
			Subtitle:   subtitle,
			IsComposed: true,
		},
	}
}

// 获取父级类型
func (this *HTTPBox) Super() *Element {
	return &this.Element
}

// 组合
func (this *HTTPBox) Compose() string {
	value := maps.NewMap(this.Value)
	url := value.GetString("url")
	method := value.GetString("method")
	timeout := value.GetString("timeout")
	textBody := value.GetString("textBody")

	if len(method) == 0 {
		method = "GET"
	}
	if len(timeout) == 0 {
		timeout = "5s"
	}

	this.Javascript = `
this.` + this.Namespace + `_httpBox_headers = ` + stringutil.JSONEncode(value.GetSlice("headers")) + `;
this.` + this.Namespace + `_httpBox_params = ` + stringutil.JSONEncode(value.GetSlice("params")) + `;
this.` + this.Namespace + `_httpBox_textBody = ` + stringutil.JSONEncode(textBody) + `;
`

	this.Attr("url", url)
	this.Attr("comment", this.Comment)
	this.Attr("method", method)
	this.Attr("timeout", timeout)
	this.Attr(":headers", this.Namespace+"_httpBox_headers")
	this.Attr(":params", this.Namespace+"_httpBox_params")
	this.Attr(":text-body", this.Namespace+"_httpBox_textBody")

	attrs := this.ComposeAttrs(this.Attrs)
	return `<tbody is="http-box" comment="` + html.EscapeString(this.Comment) + `" prefix="` + this.Namespace + `" ` + attrs + `></tbody>`
}

func (this *HTTPBox) ApplyRequest(req *http.Request) (value interface{}, skip bool, err error) {
	method := req.Form.Get(this.Namespace + "_method")
	url := req.Form.Get(this.Namespace + "_url")
	timeout := req.Form.Get(this.Namespace+"_timeout") + "s"
	textBody := req.Form.Get(this.Namespace + "_textBody")
	headers := []*shared.Variable{}
	{
		names, ok := req.Form[this.Namespace+"_headerNames"]
		if ok {
			values, ok := req.Form[this.Namespace+"_headerValues"]
			if ok {
				for index, name := range names {
					header := &shared.Variable{
						Name: name,
					}
					if index < len(values) {
						header.Value = values[index]
					}
					headers = append(headers, header)
				}
			}
		}
	}

	params := []*shared.Variable{}
	{
		names, ok := req.Form[this.Namespace+"_paramNames"]
		if ok {
			values, ok := req.Form[this.Namespace+"_paramValues"]
			if ok {
				for index, name := range names {
					param := &shared.Variable{
						Name: name,
					}
					if index < len(values) {
						param.Value = values[index]
					}
					params = append(params, param)
				}
			}
		}
	}
	return map[string]interface{}{
		"method":   method,
		"url":      url,
		"timeout":  timeout,
		"textBody": textBody,
		"headers":  headers,
		"params":   params,
	}, false, nil
}
