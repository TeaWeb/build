var widget = new widgets.Widget({
	"name": "测试Widget",
	"code": "demo@tea",
	"author": "我是测试的",
	"version": "0.0.1"
});

widget.run = function () {
	var chart = new charts.HTMLChart();
	chart.options.name = "测试";
	chart.options.columns = 2;
	chart.html = "测试HTML";
	chart.render();
};
