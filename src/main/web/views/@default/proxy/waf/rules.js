Tea.context(function () {
	this.$delay(function () {
		this.sortable();
	});

	this.enableGroup = function (groupId) {
		if (!window.confirm("确定要启用这个分组吗？")) {
			return;
		}
		this.$post("/proxy/waf/group/on")
			.params({
				"wafId": this.config.id,
				"groupId": groupId
			})
			.refresh();
	};

	this.disableGroup = function (groupId) {
		if (!window.confirm("确定要停用这个分组吗？")) {
			return;
		}
		this.$post("/proxy/waf/group/off")
			.params({
				"wafId": this.config.id,
				"groupId": groupId
			})
			.refresh();
	};

	this.deleteGroup = function (groupId) {
		if (!window.confirm("确定要删除这个分组吗？")) {
			return;
		}
		this.$post("/proxy/waf/group/delete")
			.params({
				"wafId": this.config.id,
				"groupId": groupId
			})
			.refresh();
	};

	/**
	 * 拖动排序
	 */
	this.sortable = function () {
		if (this.groups.length == 0) {
			return;
		}
		var box = this.$find("#groups-table")[0];
		var that = this;
		Sortable.create(box, {
			draggable: "tbody",
			handle: ".icon.handle",
			onStart: function () {

			},
			onUpdate: function (event) {
				var newIndex = event.newIndex;
				var oldIndex = event.oldIndex;
				that.$post("/proxy/waf/group/move")
					.params({
						"wafId": that.config.id,
						"fromIndex": oldIndex,
						"toIndex": newIndex,
						"inbound": that.inbound ? 1 : 0
					})
					.refresh();
			}
		});
	};
});