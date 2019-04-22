Tea.context(function () {
	this.$delay(function () {
		this.sortable();
	});

	this.enableRule = function (setId) {
		if (!window.confirm("确定要启用此规则集吗？")) {
			return;
		}
		this.$post("/proxy/waf/group/rule/on")
			.params({
				"wafId": this.config.id,
				"groupId": this.group.id,
				"setId": setId
			})
			.refresh();
	};

	this.disableRule = function (setId) {
		if (!window.confirm("确定要停用此规则集吗？")) {
			return;
		}
		this.$post("/proxy/waf/group/rule/off")
			.params({
				"wafId": this.config.id,
				"groupId": this.group.id,
				"setId": setId
			})
			.refresh();
	};

	this.deleteRule = function (setId) {
		if (!window.confirm("确定要删除这个规则集吗？")) {
			return;
		}
		this.$post("/proxy/waf/group/rule/delete")
			.params({
				"wafId": this.config.id,
				"groupId": this.group.id,
				"setId": setId
			})
			.refresh();
	};

	/**
	 * 拖动排序
	 */
	this.sortable = function () {
		if (this.sets.length == 0) {
			return;
		}
		var box = this.$find("#sets-table")[0];
		var that = this;
		Sortable.create(box, {
			draggable: "tbody",
			handle: ".icon.handle",
			onStart: function () {

			},
			onUpdate: function (event) {
				var newIndex = event.newIndex;
				var oldIndex = event.oldIndex;
				that.$post("/proxy/waf/group/rule/move")
					.params({
						"wafId": that.config.id,
						"groupId": that.group.id,
						"fromIndex": oldIndex,
						"toIndex": newIndex
					})
					.refresh();
			}
		});
	};
});