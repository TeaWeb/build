var widget = {
	"name": "HTTP状态码分布（今日）",
	"code": "status_stat@tea",
	"author": "TeaWeb",
	"version": "0.0.1"
};

widget.run = function () {
	var chart = new charts.PieChart();
	chart.options.name = "HTTP状态码分布<em>（今日）</em>";
	chart.options.columns = 1;

	var query = new logs.Query();
	query.from(times.now());
	query.attr("serverId", context.server.id);
	query.gt("status", 0);
	query.group("status");
	var result = query.count();

	chart.values = [];
	chart.labels = [];
	for (var key in result) {
		chart.labels.push(key);
	}

	chart.labels.sort();
	chart.labels.$each(function (k, v) {
		chart.values.push(result[v]["count"]);
	});
	if (chart.labels.length == 0) {
        chart = new charts.HTMLChart();
        chart.options.name = "HTTP状态码分布<em>（今日）</em>";
        chart.options.columns = 1;
        chart.html = "<p class='grey'><i class='ui icon pie chart'></i>暂时还没有请求。</p>";
	} else {
        if (chart.labels[0] == "200") {
            chart.colors = [colors.GREEN];
            colors.ARRAY.$each(function (k, v) {
                if (v != colors.GREEN) {
                    chart.colors.push(v);
                }
            });
        } else {
            chart.colors = colors.ARRAY;
        }
    }

	chart.render();
};