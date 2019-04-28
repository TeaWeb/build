Tea.context(function () {
	this.$delay(function () {
		this.sortable();
	});

	this.location = null;

	this.deleteLocation = function (locationId) {
		if (!window.confirm("确定要删除此路径配置吗？")) {
			return;
		}
		this.$post("/proxy/locations/delete")
			.params({
				"serverId": this.server.id,
				"locationId": locationId
			})
			.refresh();
	};

	this.moveUp = function (index) {
		this.$post("/proxy/locations/moveUp")
			.params({
				"serverId": this.server.id,
				"index": index
			});
	};

	this.moveDown = function (index) {
		this.$post("/proxy/locations/moveDown")
			.params({
				"serverId": this.server.id,
				"index": index
			});
	};

	/**
	 * 拖动排序
	 */
	this.sortable = function () {
		if (this.locations.length == 0) {
			return;
		}
		var box = this.$find("#locations-table")[0];
		var that = this;
		Sortable.create(box, {
			draggable: "tbody",
			handle: ".icon.handle",
			onStart: function () {

			},
			onUpdate: function (event) {
				var newIndex = event.newIndex;
				var oldIndex = event.oldIndex;
				that.$post("/proxy/locations/move")
					.params({
						"serverId": that.server.id,
						"fromIndex": oldIndex,
						"toIndex": newIndex
					});
			}
		});
	};
});