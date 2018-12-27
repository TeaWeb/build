var widget = new widgets.Widget({
	"name": "最近的错误日志",
	"code": "latest_error_log@tea",
	"author": "TeaWeb",
	"version": "0.0.1",
	"requirements": ["mongo"]
});

widget.run = function () {
	var chart = new charts.HTMLChart();
	chart.options.name = "最近的错误日志";
	chart.options.columns = 1;

	var query = new logs.Query();
	query.from(times.now());
	query.attr("serverId", context.server.id);
	query.gte("status", 400);
	query.limit(10);
	var errorLogs = query.desc().findAll();

	chart.html = "";
	chart.html += "<style type='text/css'> \
		.error-log-row { \
			margin-bottom: 0.6em; \
		} \
		.error-log-row .label { \
			padding: 2px; \
			margin: 0 4px; \
			display: inline-block; \
		} \
		.error-log-row .attach { \
			font-size: 0.9em; \
			color: grey; \
		}\
		</style>";
	var timestamp = times.now().unix();
	if (errorLogs.length == 0) {
		chart.html = "<p class='grey'><i class='icon history'></i>暂时还没有错误日志</p>";
	} else {
		errorLogs.$each(function (k, v) {
			var currentTime = times.unix(v.timestamp);
			chart.html += "<div class='error-log-row'><span>";
			chart.html += v.scheme + "://" + v.host + v.requestURI + "</span>";
			if (timestamp - v.timestamp < 600) {
				chart.html += "<span class='ui label tiny blue'>&lt;10m</span>";
			} else if (timestamp - v.timestamp < 3600) {
				chart.html += "<span class='ui label tiny blue'>&lt;1h</span>";
			}
			chart.html += "<span class='ui label tiny red'>" + v.status + "</span>";
			chart.html += "<div class='attach'>AT: " + currentTime.format("Y-m-d H:i:s") + "</div>";
			if (v.backendAddress.length > 0) {
				chart.html += "<div class='attach'>Backend:" + v.backendAddress + "</div>";
			}
			chart.html += "</div>";
		});
	}

	chart.render();
};
