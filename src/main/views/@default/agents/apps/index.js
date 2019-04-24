Tea.context(function () {
	this.hasSystemApp = false;

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
});