package forms

import (
	"github.com/iwind/TeaGo/types"
	"net/http"
)

type CheckBox struct {
	Element

	IsChecked bool   `yaml:"isChecked" json:"isChecked"`
	Label     string `yaml:"label" json:"label"`
}

func NewCheckBox(title string, subTitle string) *CheckBox {
	return &CheckBox{
		Element: Element{
			Title:    title,
			Subtitle: subTitle,
		},
	}
}

func (this *CheckBox) Super() *Element {
	return &this.Element
}

func (this *CheckBox) Compose() string {
	attrs := map[string]string{
		"name": this.Namespace + "_" + this.Code,
	}

	if this.IsChecked || types.Bool(this.Value) {
		attrs["checked"] = "checked"
	}

	attrs["value"] = "1"

	return `
<div class="ui checkbox">
<input type="checkbox"` + this.ComposeAttrs(attrs) + `/>
<label>` + this.Label + `</label>
</div>
`
}

func (this *CheckBox) ApplyRequest(req *http.Request) (value interface{}, skip bool, err error) {
	value = req.Form.Get(this.Namespace + "_" + this.Code)
	return types.Bool(value), false, nil
}
