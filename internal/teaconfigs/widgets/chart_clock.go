package widgets

import "github.com/iwind/TeaGo/utils/string"

// 时钟
type ClockChart struct {
}

func (this *ClockChart) AsJavascript(options map[string]interface{}) (code string, err error) {
	return `
var chart = new charts.Clock();
chart.options = ` + stringutil.JSONEncode(options) + `;

var latest = NewQuery().latest(1);
if (latest.length > 0) {
	chart.timestamp = parseInt(new Date().getTime() / 1000) - (latest[0].createdAt - latest[0].value.timestamp);
}
chart.render();
`, nil
}
