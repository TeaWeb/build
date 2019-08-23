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
				"dayTo": this.chart.dayTo
			})
			.success(function (resp) {
				// output
				resp.data.output.$each(function (k, v) {
					console.log("[widget]" + v);
				});

				// charts
				this.charts = resp.data.charts;
				new ChartRender(this.charts);
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
});