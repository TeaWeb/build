package widgets

import "github.com/iwind/TeaGo/utils/string"

// HTML Chart
type HTMLChart struct {
	HTML string `yaml:"html" json:"html"`
}

func (this *HTMLChart) AsJavascript(options map[string]interface{}) (code string, err error) {
	return `
var chart = new charts.HTMLChart();
chart.options = ` + stringutil.JSONEncode(options) + `;
chart.html = ` + stringutil.JSONEncode(this.HTML) + `;
chart.render();
`, nil
}
