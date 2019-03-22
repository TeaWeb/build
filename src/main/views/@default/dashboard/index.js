Tea.context(function () {
	this.isLoaded = false;
	this.logsVisible = false;
	this.charts = [];
	this.logs = [];
	this.agentId = "local";

	this.$delay(function () {
		this.testMongo();
		this.loadData();
		this.loadLogs();
	});

	this.loadData = function () {
		this.$post("/agents/board")
			.params({
				"agentId": "local"
			})
			.success(function (resp) {
				this.charts = resp.data.charts;
				new ChartRender(this.charts);
			})
			.done(function () {
				this.$delay(function () {
					this.isLoaded = true;
				});
				this.$delay(function () {
					this.loadData();
				}, 5000);
			});
	};

	this.loadLogs = function () {
		this.$get("/dashboard/logs")
			.params({})
			.success(function (resp) {
				this.logs = resp.data.logs.$map(function (k, v) {
					v.requestTime = Math.floor(v.requestTime * 1000000) / 1000;
					return v;
				});

				this.renderQPS(resp.data.qps);

				this.logsVisible = true;
			})
			.done(function () {
				this.$delay(function () {
					this.loadLogs();
				}, 3000);
			});
	};

	this.removeChart = function (appId, itemId, chartId) {
		this.$post("/agents/board/removeChart")
			.params({
				"agentId": this.agentId,
				"appId": appId,
				"itemId": itemId,
				"chartId": chartId
			})
			.refresh();

		return false;
	};

	var lastQPSMax = 0;
	this.renderQPS = function (qps) {
		var chartBox = document.getElementById("qps-chart-box");
		var chart = echarts.init(chartBox);
		var max = 10000;
		if (qps < 100) {
			max = 100;
		} else if (qps < 1000) {
			max = 1000;
		} else if (qps < 10000) {
			max = 10000;
		} else {
			max = 100000;
		}
		if (max < lastQPSMax) {
			max = lastQPSMax;
		} else {
			lastQPSMax = max;
		}
		var options = {
			"name": "",
			"min": 0,
			"max": max,
			"detail": "QPS",
			"value": qps,
			"unit": "Req"
		};
		var option = {
			textStyle: {
				fontFamily: "Lato,'Helvetica Neue',Arial,Helvetica,sans-serif"
			},
			title: {
				text: options.name,
				top: 1,
				bottom: 0,
				x: "center",
				textStyle: {
					fontSize: 12,
					fontWeight: "bold",
					fontFamily: "Lato,'Helvetica Neue',Arial,Helvetica,sans-serif"
				}
			},
			legend: {
				data: [""]
			},
			xAxis: {
				data: []
			},
			yAxis: {},
			series: [{
				name: '',
				type: 'gauge',
				min: options.min,
				max: options.max,

				data: [
					{
						"name": options.detail,
						"value": Math.round(options.value * 100) / 100
					}
				],
				radius: "80%",
				center: ["50%", (options.name != null && options.name.length > 0) ? "60%" : "50%"],

				splitNumber: 5,
				splitLine: {
					length: 6
				},

				axisLine: {
					lineStyle: {
						width: 8
					}
				},
				axisTick: {
					show: true
				},
				axisLabel: {
					formatter: function (v) {
						return v;
					},
					textStyle: {
						fontSize: 8
					}
				},
				detail: {
					formatter: function (v) {
						return v + options.unit;
					},
					textStyle: {
						fontSize: 12
					}
				},

				pointer: {
					width: 2
				}
			}],

			grid: {
				left: -2,
				right: 0,
				bottom: 0,
				top: 0
			},
			axisPointer: {
				show: false
			},
			tooltip: {
				formatter: 'X:{b0} Y:{c0}',
				show: false
			},
			animation: true
		};

		chart.setOption(option);
	};
});