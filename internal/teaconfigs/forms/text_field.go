package forms

import (
	"fmt"
	"github.com/iwind/TeaGo/types"
	"net/http"
)

type TextField struct {
	Element     `yaml:",inline"`
	MaxLength   int    `yaml:"maxLength" json:"maxLength"`
	Placeholder string `yaml:"placeholder" json:"placeholder"`
	Size        int    `yaml:"size" json:"size"`
	RightLabel  string `yaml:"rightLabel" json:"rightLabel"`
}

func NewTextField(title string, subTitle string) *TextField {
	return &TextField{
		Element: Element{
			Title:    title,
			Subtitle: subTitle,
		},
	}
}

func (this *TextField) Super() *Element {
	return &this.Element
}

func (this *TextField) Compose() string {
	attrs := map[string]string{}
	if this.MaxLength > 0 {
		attrs["maxlength"] = fmt.Sprintf("%d", this.MaxLength)
	}
	if len(this.Placeholder) > 0 {
		attrs["placeholder"] = this.Placeholder
	}
	attrs["value"] = types.String(this.Value)
	if this.Size > 0 {
		attrs["size"] = fmt.Sprintf("%d", this.Size)
	}
	attrs["name"] = this.Namespace + "_" + this.Code

	if len(this.RightLabel) == 0 {
		return `<input type="text" ` + this.ComposeAttrs(attrs) + ` />`
	} else {
		return `<div class="ui input right labeled"><input type="text" ` + this.ComposeAttrs(attrs) + `> <label class="ui label">` + this.RightLabel + `</label></div>`
	}
}

func (this *TextField) ApplyRequest(req *http.Request) (value interface{}, skip bool, err error) {
	value = req.Form.Get(this.Namespace + "_" + this.Code)
	return value, false, nil
}
