package forms

import (
	"fmt"
	"github.com/iwind/TeaGo/types"
	"net/http"
)

type TextBox struct {
	Element
	MaxLength   int
	Placeholder string
	Cols        int
	Rows        int
}

func NewTextBox(title string, subTitle string) *TextBox {
	return &TextBox{
		Element: Element{
			Title:    title,
			Subtitle: subTitle,
		},
	}
}

func (this *TextBox) Super() *Element {
	return &this.Element
}

func (this *TextBox) Compose() string {
	attrs := map[string]string{}
	if this.MaxLength > 0 {
		attrs["maxlength"] = fmt.Sprintf("%d", this.MaxLength)
	}

	if len(this.Placeholder) > 0 {
		attrs["placeholder"] = this.Placeholder
	}

	if this.Cols > 0 {
		attrs["cols"] = fmt.Sprintf("%d", this.Cols)
	}

	if this.Rows > 0 {
		attrs["rows"] = fmt.Sprintf("%d", this.Rows)
	}

	valueString := types.String(this.Value)
	attrs["name"] = this.Namespace + "_" + this.Code

	return `<textarea ` + this.ComposeAttrs(attrs) + `>` + valueString + `</textarea>`
}

func (this *TextBox) ApplyRequest(req *http.Request) (value interface{}, skip bool, err error) {
	value = req.Form.Get(this.Namespace + "_" + this.Code)
	return value, false, nil
}
