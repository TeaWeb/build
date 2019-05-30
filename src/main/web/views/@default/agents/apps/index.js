Tea.context(function () {
	this.hasSystemApp = false;

	this.$delay(function () {
		this.sortable();
	});

	if (this.apps != null) {
		this.hasSystemApp = this.apps.$exist(function (k, v) {
			return v.id == "system";
		});
	}

	this.deleteApp = function (appId) {
		if (!window.confirm("确定要删除此App吗？")) {
			return;
		}
		this.$post("/agents/apps/delete")
			.params({
				"agentId": this.agentId,
				"appId": appId
			})
			.refresh();
	};

	this.addSystemApp = function () {
		if (!window.confirm("确定要添加内置的系统App吗？")) {
			return;
		}
		this.$post("/agents/board/initDefaultApp")
			.params({
				"agentId": this.agentId
			})
			.refresh();
	};

	/**
	 * 拖动排序
	 */
	this.sortable = function () {
		if (this.apps.length == 0) {
			return;
		}
		var box = this.$find("#apps-table")[0];
		var that = this;
		Sortable.create(box, {
			draggable: "tbody",
			handle: ".icon.handle",
			onStart: function () {

			},
			onUpdate: function (event) {
				var newIndex = event.newIndex;
				var oldIndex = event.oldIndex;
				that.$post("/agents/apps/move")
					.params({
						"agentId": that.agentId,
						"fromIndex": oldIndex,
						"toIndex": newIndex
					});
			}
		});
	};
});