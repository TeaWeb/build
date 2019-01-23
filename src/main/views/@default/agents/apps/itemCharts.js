Tea.context(function () {
	this.from = encodeURIComponent(window.location.toString());
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
				"agentId": this.agentId,
				"appId": this.app.id,
				"itemId": this.item.id,
				"from": this.from
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
				this.isLoaded = true;
				this.$delay(function () {
					this.loadCharts();
				}, (this.intervalSeconds > 0) ? this.intervalSeconds * 1000 : 10 * 1000);
			});
	};

	this.deleteChart = function (chartId) {
		if (!window.confirm("确定要删除这个图表吗？")) {
			return false;
		}
		this.$post("/agents/apps/deleteItemChart")
			.params({
				"agentId": this.agentId,
				"appId": this.app.id,
				"itemId": this.item.id,
				"chartId": chartId
			})
			.refresh();
		return false;
	};
});