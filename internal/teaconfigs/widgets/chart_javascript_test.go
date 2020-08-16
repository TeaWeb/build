package widgets

import (
	"github.com/iwind/TeaGo/maps"
	"testing"
	"time"
)

func TestJavascriptChart_AsJavascript(t *testing.T) {
	c := new(JavascriptChart)
	c.Code = `
var chart = new charts.HTMLChart();
chart.options.name = "Hello";
chart.render();
`
	before := time.Now()
	t.Log(c.AsJavascript(maps.Map{
		"name":    "Hello,World",
		"columns": 2,
	}))
	t.Log(time.Since(before).Seconds(), "seconds")
}
