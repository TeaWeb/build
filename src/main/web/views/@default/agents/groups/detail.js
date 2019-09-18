Tea.context(function () {
	this.keyword = "";
	this.filterOnline = "-1";
	var allAgents = this.agents;
	this.hasAgents = allAgents.length > 0;
	this.countAllAgents = allAgents.length;

	this.$delay(function () {
		this.sortable();
	});

	/**
	 * 密钥
	 */
	this.generateKey = function () {
		if (this.group.key != null && this.group.key.length > 0) {
			if (!window.confirm("确定要重新生成密钥吗？")) {
				return;
			}
		}
		this.$post(".generateKey")
			.params({
				"groupId": this.group.id
			})
			.refresh();
	};

	/**
	 * 生成
	 */
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
			v.isChecked = false;
			return true;
		});
		this.selectedAgents = [];
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

	/**
	 * 全选删除
	 */
	this.selectedAgents = [];

	this.selectAllAgents = function () {
		this.agents = this.agents.$map(function (k, v) {
			v.isChecked = true;
			return v;
		});
		this.selectedAgents = this.agents;
	};

	this.changeAgentChecked = function () {
		this.$delay(function () {
			this.selectedAgents = this.agents.$findAll(function (k, v) {
				return v.isChecked;
			});
		});
	};

	this.isDeleting = false;

	this.deleteAgents = function () {
		if (this.selectedAgents.length == 0) {
			alert("请选择要删除的主机");
			return;
		}

		if (!window.confirm("确定要删除选中的主机吗？")) {
			return;
		}

		this.isDeleting = true;

		this.$post("/agents/deleteAgents")
			.params({
				"agentIds": this.selectedAgents.$map(function (k, v) {
					return v.id
				})
			})
			.refresh();
	};

	/**
	 * 密钥
	 */
	this.selectedKeyIndex = -1;
	this.keysVisible = false;

	this.selectKey = function (index) {
		this.selectedKeyIndex = index;
	};

	this.showKeys = function () {
		this.keysVisible = !this.keysVisible;
	};
});