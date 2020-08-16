package scripts

import (
	"testing"
)

func TestEngine_Run(t *testing.T) {
	engine := NewEngine()
	err := engine.RunCode(`var widget = new widgets.Widget({
	"name": "测试Widget",
	"code": "test_stat@tea",
	"author": "我是测试的",
	"version": "0.0.1"
});

widget.run = function () {
	var chart = new charts.HTMLChart();
	chart.html = "测试HTML Chart";
	chart.render();
};`)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(engine.Charts())
}

func TestEngine_Cache(t *testing.T) {
	engine := NewEngine()
	err := engine.RunCode(`
var widget = new widgets.Widget({});
widget.run = function () {
	caches.set("a", "b");
	console.log(caches.get("a"));
};
`)
	if err != nil {
		t.Fatal(err)
	}
}
