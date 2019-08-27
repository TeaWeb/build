Tea.context(function () {
	this.$delay(function () {
		var that = this;
		this.$find("input[name='name']").focus();
		teaweb.datepicker(document.getElementById("day-from-input"), function (day) {
			that.chart.dayFrom = day;
			that.loadChart();
		});
		teaweb.datepicker(document.getElementById("day-to-input"), function (day) {
			that.chart.dayTo = day;
			that.loadChart();
		});

		this.loadChart();
	});

	/**
	 * 显示类型
	 */
	this.tab = "chart";
	this.exportType = "";
	this.exportArray = [];
	this.exportTitles = [];

	this.selectTab = function (tab) {
		this.tab = tab;

		if (tab == "data") {
			this.exportType = "data";
		} else if (tab == "chart") {
			this.exportType = "";
		}

		if (tab == "chart" || tab == "data") {
			this.loadChart();
		}
	};

	/**
	 * 导出CSV
	 */
	this.exportCSV = function () {
		this.exportType = "csv";
		var params = {
			"name": this.chart.name,
			"agentId": this.agentId,
			"appId": this.appId,
			"itemId": this.itemId,
			"chartId": this.chartId,
			"timeType": this.chart.timeType,
			"timePast": this.chart.timePast,
			"dayFrom": this.chart.dayFrom,
			"dayTo": this.chart.dayTo,
			"export": this.exportType
		};

		window.location = "/agents/board/exportChartData?" + Tea.serialize(params);
	};

	/**
	 * 图表
	 */
	this.charts = [];

	this.changeName = function () {
		this.loadChart();
	};

	this.changeTimeType = function () {
		this.loadChart();
	};

	this.changeTimePast = function () {
		this.loadChart();
	};

	this.isLoading = false;
	this.loadChart = function () {
		this.isLoading = true;
		this.successVisible = false;

		this.exportArray = [];
		this.exportTitles = [];

		this.$post("$")
			.params({
				"name": this.chart.name,
				"agentId": this.agentId,
				"appId": this.appId,
				"itemId": this.itemId,
				"chartId": this.chartId,
				"timeType": this.chart.timeType,
				"timePast": this.chart.timePast,
				"dayFrom": this.chart.dayFrom,
				"dayTo": this.chart.dayTo,
				"export": this.exportType
			})
			.success(function (resp) {
				// output
				resp.data.output.$each(function (k, v) {
					console.log("[widget]" + v);
				});

				// charts
				if (this.exportType.length == 0) {
					this.charts = resp.data.charts;
					new ChartRender(this.charts);
				} else if (this.exportType == "data") {
					var result = resp.data.result;
					if (result == null) {
						return;
					}
					if (result instanceof Array) {
						var that = this;
						this.exportArray = result.$map(function (k, v) {
							var time = "";
							switch (resp.data.timeUnit) {
								case "":
									time = v.timeFormat.minute;
									break;
								case "SECOND":
									time = v.timeFormat.second;
									break;
								case "MINUTE":
									time = v.timeFormat.minute;
									break;
								case "HOUR":
									time = v.timeFormat.hour;
									break;
								case "DAY":
									time = v.timeFormat.day;
									break;
								case "MONTH":
									time = v.timeFormat.month;
									break;
								case "YEAR":
									time = v.timeFormat.year;
									break;
							}

							// 展开value
							if (v.value != null && typeof (v.value) == "object") {
								that.extractTitles("", v.value);
							}

							return {
								"value": v.value,
								"time": that.formatTime(time)
							};
						});
					}
				}
			})
			.done(function () {
				this.isLoading = false;
			});
	};

	this.successVisible = false;
	this.submit = function () {
		var form = new FormData(document.getElementById("update-chart-form"))
		this.$post(".updateChart")
			.params(form)
			.success(function () {
				this.successVisible = true;
			});
	};

	this.formatTime = function (t) {
		if (t.length <= 4) { // year
			return t;
		}
		if (t.length == 6) { // month
			return t.substr(0, 4) + "-" + t.substr(4, 2);
		}
		if (t.length == 8) { // day
			return this.formatTime(t.substr(0, 6)) + "-" + t.substr(6, 2);
		}
		if (t.length == 10) { // hour
			return this.formatTime(t.substr(0, 8)) + " " + t.substr(8, 2)
		}
		if (t.length == 12) { // minute
			return this.formatTime(t.substr(0, 10)) + ":" + t.substr(10, 2);
		}
		if (t.length == 14) { // second
			return this.formatTime(t.substr(0, 12)) + ":" + t.substr(12, 2);
		}
		return t;
	};

	this.getValue = function (v, k) {
		if (v == null) {
			return null;
		}
		var pieces = k.split(".");
		if (pieces.length == 1) {
			return v[pieces[0]];
		}
		if (typeof (v) != "object") {
			return v;
		}
		return this.getValue(v[pieces[0]], pieces.slice(1).join("."));
	};

	this.extractTitles = function (prefix, v) {
		if (v == null) {
			return;
		}
		for (var k in v) {
			if (!v.hasOwnProperty(k)) {
				continue;
			}
			var v1 = v[k];
			var key = k;
			if (prefix != null && prefix.length > 0) {
				key = prefix + "." + key;
			}
			if (v1 != null && typeof (v1) == "object" && !(v1 instanceof Array)) {
				this.extractTitles(key, v1);
				continue;
			}

			if (!this.exportTitles.$contains(key)) {
				this.exportTitles.push(key);
			}
		}
	};
});