Tea.context(function () {
	this.$delay(function () {
		this.sortable()
	});

	this.useChart = function (chart, serverId, widgetId, chartId) {
		this.$post("/proxy/board/useChart")
			.params({
				"serverId": serverId,
				"widgetId": widgetId,
				"chartId": chartId,
				"type": this.boardType
			})
			.success(function () {
				this.usingCharts.push({
					"id": chart.id,
					"name": chart.name,
					"widgetId": widgetId,
					"on": chart.on,
					"columns": chart.columns,
					"widget": {
						"id": widgetId
					}
				});

				chart.isUsing = true;
			})
	};

	this.cancelChart = function (chart, serverId, widgetId, chartId) {
		this.$post("/proxy/board/cancelChart")
			.params({
				"serverId": serverId,
				"widgetId": widgetId,
				"chartId": chartId,
				"type": this.boardType
			})
			.success(function () {
				var widget = this.widgets.$find(function (k, v) {
					return widgetId == v.id;
				});
				var c = widget.charts.$find(function (k, v) {
					return v.id == chartId
				});
				c.isUsing = false;

				this.usingCharts.$removeValue(chart);
			});
	};

	/**
	 * 排序
	 */
	this.moveSuccess = false;

	this.sortable = function () {
		var box = this.$find(".using-charts-box")[0];
		var that = this;
		Sortable.create(box, {
			draggable: ".chart-box",
			onStart: function () {

			},
			onUpdate: function (event) {
				var newIndex = event.newIndex;
				var oldIndex = event.oldIndex;

				that.$post("/proxy/board/moveChart")
					.params({
						"serverId": that.server.id,
						"type": that.boardType,
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