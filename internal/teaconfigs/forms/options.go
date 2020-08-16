package forms

import (
	"github.com/iwind/TeaGo/types"
	"html"
	"net/http"
)

type Options struct {
	Element `yaml:",inline"`

	Options []*Option
}

type Option struct {
	Text  string
	Value string
}

func NewOptions(title string, subtitle string) *Options {
	return &Options{
		Element: Element{
			Title:    title,
			Subtitle: subtitle,
		},
	}
}

func (this *Options) AddOption(text string, value string) {
	this.Options = append(this.Options, &Option{
		Text:  text,
		Value: value,
	})
}

func (this *Options) Super() *Element {
	return &this.Element
}

func (this *Options) Compose() string {
	result := `<select ` + this.ComposeAttrs(map[string]string{
		"name": this.Namespace + "_" + this.Code,
	}) + ` class="ui dropdown">`
	for _, o := range this.Options {
		result += `<option value="` + html.EscapeString(o.Value) + `"`
		if o.Value == types.String(this.Value) {
			result += ` selected="selected"`
		}
		result += ">" + html.EscapeString(o.Text) + `</option>`
	}
	result += "</select>"
	return result
}

func (this *Options) ApplyRequest(req *http.Request) (value interface{}, skip bool, err error) {
	value = req.Form.Get(this.Namespace + "_" + this.Code)
	return value, false, nil
}
