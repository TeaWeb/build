Tea.context(function () {
	this.isLoaded = false;

	this.$delay(function () {
		this.loadCharts();
	});

	/**
	 * 加载图表
	 */
	this.charts = [];

	this.loadCharts = function () {
		this.$post("$")
			.params({
				"agentId": this.agentId
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
				this.$delay(function () {
					this.isLoaded = true;
				});
				this.$delay(function () {
					this.loadCharts();
				}, (this.intervalSeconds > 0) ? this.intervalSeconds * 1000 : 10 * 1000);
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
});