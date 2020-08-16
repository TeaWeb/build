package widgets

import "github.com/iwind/TeaGo/utils/string"

// URL Chart
type URLChart struct {
	URL string `yaml:"url" json:"url"`
}

func (this *URLChart) AsJavascript(options map[string]interface{}) (code string, err error) {
	return `
var chart = new charts.URLChart();
chart.options = ` + stringutil.JSONEncode(options) + `;
chart.url = ` + stringutil.JSONEncode(this.URL) + `;
chart.render();
`, nil
}
