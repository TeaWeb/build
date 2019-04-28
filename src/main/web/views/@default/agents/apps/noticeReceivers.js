Tea.context(function () {
	this.deleteReceiver = function (level, receiverId) {
		if (!window.confirm("确定要删除此接收人吗？")) {
			return;
		}
		this.$post("/agents/apps/deleteNoticeReceivers")
			.params({
				"agentId": this.agentId,
				"appId": this.app.id,
				"level": level,
				"receiverId": receiverId
			})
			.refresh();
	};
});