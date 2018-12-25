var widget = {
	"name": "请求耗时",
	"code": "request_time@tea",
	"author": "TeaWeb",
	"version": "0.0.1"
};

widget.run = function () {
	var chart = new charts.LineChart();
	chart.options.name = "请求耗时趋势<em>（最近五分钟平均时间，单位：ms）</em>";
	chart.options.columns = 2;
	chart.xShowTick = false;

	var timeList = [];
	var now = times.now();
	var passedTimestamp = times.now().addTime(0, 0, 0, 0, -5, 0).unix();
	passedTimestamp -= passedTimestamp % 30;
	var passed = times.unix(passedTimestamp);
	while (true) {
		var passedTo = passed.addTime(0, 0, 0, 0, 0, 30);
		timeList.push(passed.format("YmdHis") + "_" + passedTo.format("YmdHis"));
		passed = passedTo;
		if (passed.unix() > now.unix()) {
			timeList.push(passed.format("YmdHis") + "_" + passedTo.format("YmdHis"));
			break;
		}
	}

	var values = [];
	timeList.$each(function (k, v) {
		var pieces = v.split("_");
		var timeFrom = pieces[0];
		var timeTo = pieces[1];
		var query = new logs.Query();
		query.attr("serverId", context.server.id);
		query.from(now);
		query.to(now);
		query.gte("timeFormat.second", timeFrom);
		query.lte("timeFormat.second", timeTo);
		var costMs = query.avg("requestTime") * 1000;
		values.push(costMs);
	});

	chart.labels = [];
	chart.labels.$fill("", values.length);

	var line1 = new charts.Line();
	line1.isFilled = true;
	line1.values = values;
	if (values.length > 0) {
		var avg = values.$sum(function (k, v) {
			return v;
		}) / values.length;
		if (avg > 500) {
			line1.color = colors.RED;
		}
	}

	var max = values.$max();
	if (max < 1) {
		chart.max = 1;
	} else if (max < 10) {
		chart.max = 10;
	} else if (max < 100) {
		chart.max = 100;
	}
	chart.addLine(line1);

	chart.render();
};