var widget = {
	"name": "实时请求",
	"code": "request_realtime@tea",
	"author": "TeaWeb",
	"version": "0.0.1"
};

widget.run = function () {
	var chart = new charts.LineChart();
	chart.options.name = "实时请求<em>（Req/s）</em>";
	chart.options.columns = 2;
	chart.xShowTick = false;

	var timeList = [];
	var now = times.now();
	var passedTimestamp = times.now().addTime(0, 0, 0, 0, -1, 0).unix();
	passedTimestamp -= passedTimestamp % 4;
	var passed = times.unix(passedTimestamp);
	while (true) {
		timeList.push(passed.format("YmdHis"));
		passed = passed.addTime(0, 0, 0, 0, 0, 4);
		if (passed.unix() > now.unix()) {
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
		var count = query.count();
		values.push(count);
	});

	chart.labels = [];
	chart.labels.$fill("", timeList.length);

	var line1 = new charts.Line();
	line1.isFilled = true;
	line1.values = values;
	chart.addLine(line1);

	chart.render();
};