Tea.context(function () {
	this.keyword = "";
	this.filterOnline = "-1";
	var allAgents = this.agents;
	this.hasAgents = allAgents.length > 0;

	this.$delay(function () {
		this.sortable();
	});

	this.changeKeyword = function () {
		this.filter();
	};

	this.resetKeyword = function () {
		this.keyword = "";
		this.filterOnline = "-1";
		this.$delay(function () {
			this.changeKeyword();
		});
	};

	this.changeOnlineFilter = function () {
		this.filter();
	};

	this.filter = function () {
		if (this.keyword.length == 0 && this.filterOnline == "-1") {
			this.agents = allAgents;
			return;
		}
		var keyword = this.keyword;
		var filterOnline = this.filterOnline;
		this.agents = allAgents.$filter(function (k, v) {
			if (keyword.length > 0) {
				if (!teaweb.match(v.name + " " + v.host, keyword)) {
					return false;
				}
			}
			if (filterOnline != "-1") {
				if (filterOnline == "0" && v.isWaiting) {
					return false;
				}
				if (filterOnline == "1" && !v.isWaiting) {
					return false;
				}
			}
			return true;
		});
	};

	/**
	 * 拖动排序
	 */
	this.sortable = function () {
		if (this.agents.length == 0) {
			return;
		}
		var box = this.$find("#agents-table")[0];
		var that = this;
		Sortable.create(box, {
			draggable: "tbody",
			handle: ".icon.handle",
			onStart: function () {

			},
			onUpdate: function (event) {
				var newIndex = event.newIndex;
				var oldIndex = event.oldIndex;
				var toId = allAgents[newIndex].id;
				var fromId = allAgents[oldIndex].id;
				that.$post("/agents/move")
					.params({
						"fromId": fromId,
						"toId": toId
					})
					.refresh();
			}
		});
	};
});