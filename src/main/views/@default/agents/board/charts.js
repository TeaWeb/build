Tea.context(function () {
	this.$delay(function () {
		this.sortable();
	});

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

	/**
	 * 排序
	 */
	this.moveSuccess = false;

	this.sortable = function () {
		var box = this.$find(".chart-box")[0];
		var that = this;
		Sortable.create(box, {
			draggable: "span",
			onStart: function () {

			},
			onUpdate: function (event) {
				var newIndex = event.newIndex;
				var oldIndex = event.oldIndex;

				that.$post("/agents/board/moveChart")
					.params({
						"agentId": that.agentId,
						"oldIndex": oldIndex,
						"newIndex": newIndex
					})
					.success(function () {
						that.moveSuccess = true;
						this.$delay(function () {
							window.location.reload();
						}, 1000);
					});
			}
		});
	};
});