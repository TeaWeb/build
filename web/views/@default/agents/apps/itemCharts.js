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
	this.hasWrongCharts = false;

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

				// 和预计的数量是否一致
				this.hasWrongCharts = (this.charts.length != this.item.charts.length);
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

	this.removeChartFormBoard = function (chartId) {
		if (!window.confirm("确定从看板中移除这个图表吗？")) {
			return false;
		}
		this.$post("/agents/board/removeChart")
			.params({
				"agentId": this.agentId,
				"appId": this.app.id,
				"itemId": this.item.id,
				"chartId": chartId
			})
			.refresh();
		return false;
	};

	this.addChartToBoard = function (chartId) {
		if (!window.confirm("确定把这个图表添加到看板吗？")) {
			return false;
		}
		this.$post("/agents/board/addChart")
			.params({
				"agentId": this.agentId,
				"appId": this.app.id,
				"itemId": this.item.id,
				"chartId": chartId
			})
			.refresh();
		return false;
	};

	/**
	 * 添加数据源内置图表
	 */
	this.addDefaultCharts = function () {
		if (!window.confirm("确定要添加数据源内置图表吗？")) {
			return;
		}
		this.$post("/agents/apps/addDefaultCharts")
			.params({
				"agentId": this.agentId,
				"appId": this.app.id,
				"itemId": this.item.id
			})
			.refresh();
	};
});