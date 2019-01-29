Tea.context(function () {
	this.setRead = function (scope, noticeIds, msg) {
		if (msg != null) {
			if (!window.confirm(msg)) {
				return;
			}
		}
		this.$post("/agents/notices/setRead")
			.params({
				"scope": scope,
				"agentId": this.agentId,
				"noticeIds": (noticeIds != null) ? noticeIds : this.notices.$map(function (k, v) {
					return v.id;
				})
			})
			.success(function () {
				if (scope == "page") {
					window.location.reload();
				} else {
					window.location = "/agents/notices?agentId=" + this.agentId;
				}
			});
	};
});