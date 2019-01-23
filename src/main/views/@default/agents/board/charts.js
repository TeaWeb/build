Tea.context(function () {
	this.addChart = function (chart) {
		this.$post("/agents/board/addChart")
			.params({
				"agentId": this.agentId,
				"appId": chart.app.id,
				"itemId": chart.item.id,
				"chartId": chart.id
			})
			.success(function () {
				chart.isUsing = true;
			});
	};

	this.removeChart = function (chart) {
		this.$post("/agents/board/removeChart")
			.params({
				"agentId": this.agentId,
				"appId": chart.app.id,
				"itemId": chart.item.id,
				"chartId": chart.id
			})
			.success(function () {
				chart.isUsing = false;
			});
	};

});