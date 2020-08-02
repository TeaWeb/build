Tea.context(function () {
	this.isLoaded = false;

	this.$delay(function () {
		this.loadCharts();
		this.agentsSortable();
	});

	/**
	 * 加载图表
	 */
	this.charts = [];
	this.error = "";

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

				this.error = resp.data.error;
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

	/**
	 * 左侧菜单排序
	 */
	this.agentsSortable = function () {
		var that = this;
		this.$find(".sub-menu .menus-box div").each(function (k, box) {
			var items = Tea.element(box).find("a.item.sortable");
			Sortable.create(box, {
				draggable: "a.item.sortable",
				onStart: function () {

				},
				onUpdate: function (event) {
					var newIndex = event.newIndex;
					var oldIndex = event.oldIndex;

					var fromId = Tea.element(items[oldIndex]).attr("data-id");
					var toId = Tea.element(items[newIndex]).attr("data-id");

					that.$post("/agents/move")
						.params({
							"fromId": fromId,
							"toId": toId
						})
						.success(function () {
							this.$get("/agents/menu")
								.params({"agentId": this.agentId})
								.success(function (resp) {
									this.teaSubMenus.menus = [];
									this.$delay(function () {
										this.teaSubMenus = resp.data.teaSubMenus;
										this.teaSubMenus.menus.$each(function (k, menu) {
											menu.items.$each(function (k, item) {
												if (item.id == fromId && !menu.isActive) {
													that.showSubMenu(menu);
												}
											});
										});
										this.$delay(function () {
											this.agentsSortable();
										});
									});
								});
						});
				}
			});
		});
	};

	/**
	 * 重新初始化
	 */
	this.initDefaultApps = function () {
		if (!window.confirm("确定要初始化内置的系统App吗？")) {
			return;
		}
		this.$post(".initDefaultApp")
			.params({
				"agentId": this.agentId
			})
			.refresh();
	};

	/**
	 * 显示设置
	 */
	this.updateChartSetting = function (appId, itemId, chartId) {
		this.showModal("chart-setting-modal");
	};
});