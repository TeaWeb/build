Tea.context(function () {
	var that = this;
	this.from = encodeURIComponent(window.location.toString());
	this.items = [];

	this.$delay(function () {
		this.loadItems();
	});

	this.loadItems = function () {
		this.$post("$")
			.params({
				"agentId": this.agentId,
				"appId": this.app.id
			})
			.success(function (resp) {
				this.items = resp.data.items;
				this.items.$each(function (k, item) {
					item.costMs = Math.ceil(item.costMs * 1000) / 1000;

					// 阈值
					if (item.thresholds != null) {
						item.thresholds.$each(function (k, v) {
							v.levelName = that.noticeLevels.$find(function (k, v1) {
								return v.noticeLevel == v1.code;
							}).name;
						});
					}
				});

				this.$delay(function () {
					this.loadItems();
				}, 5000);
			})
	};

	this.deleteItem = function (itemId) {
		if (!window.confirm("确定要删除此监控项吗？")) {
			return;
		}
		this.$post("/agents/apps/deleteItem")
			.params({
				"agentId": this.agentId,
				"appId": this.app.id,
				"itemId": itemId
			})
			.refresh();
	};
});