var widget = new widgets.Widget({
	"name": "实时带宽",
	"code": "bandwidth_realtime@tea",
	"author": "TeaWeb",
	"version": "0.0.1",
	"requirements": ["mongo"]
});

widget.run = function () {
	var chart = new charts.LineChart();
	chart.options.name = "实时带宽<em>（KB/s）</em>";
	chart.options.columns = 2;
	chart.xShowTick = false;

	var timeList = [];
	var now = times.now();
	var passedTimestamp = times.now().addTime(0, 0, 0, 0, -1, 0).unix();
	passedTimestamp -= passedTimestamp % 2;
	var passed = times.unix(passedTimestamp);
	while (true) {
		timeList.push(passed.format("YmdHis"));
		passed = passed.addTime(0, 0, 0, 0, 0, 2);
		if (passed.unix() >= now.unix()) {
			break;
		}
	}

	var values = [];
	timeList.$each(function (k, v) {
		var query = new logs.Query();
		query.attr("serverId", context.server.id);
		query.from(now);
		query.to(now);
		query.attr("timeFormat.second", v);
		var bytesSent = query.sum("bytesSent") / 1024;
		values.push(bytesSent);
	});

	var max = values.$max();
	if (max < 1) {
        chart.options.name = "实时带宽<em>（B/s）</em>";
        values = values.$map(function (k, v) {
            return v  * 1024;
        });
    } else if (max < 10) {
		chart.max = 10;
	} else if (max > 1000) {
		chart.options.name = "实时带宽<em>（MB/s）</em>";
		values = values.$map(function (k, v) {
			return v / 1024;
		});
	}

	chart.labels = [];
	chart.labels.$fill("", timeList.length);

	var line1 = new charts.Line();
	line1.isFilled = true;
	line1.values = values;
    line1.smooth = true;
	chart.addLine(line1);

	chart.render();
};