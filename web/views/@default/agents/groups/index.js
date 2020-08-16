Tea.context(function () {
	this.$delay(function () {
		this.sortable();
	});

	this.deleteGroup = function (groupId) {
		if (!window.confirm("确定要删除此分组吗？")) {
			return;
		}
		this.$post("/agents/groups/delete")
			.params({
				"groupId": groupId
			})
			.refresh();
	};

	/**
	 * 排序
	 */
	this.sortable = function () {
		var that = this;
		this.$find("#groups-table").each(function (k, box) {
			Sortable.create(box, {
				draggable: "tbody.sortable",
				onStart: function () {

				},
				onUpdate: function (event) {
					var newIndex = event.newIndex;
					var oldIndex = event.oldIndex;

					that.$post("/agents/groups/move")
						.params({
							"fromIndex": oldIndex,
							"toIndex": newIndex
						})
						.refresh();
				}
			});
		});
	};
});