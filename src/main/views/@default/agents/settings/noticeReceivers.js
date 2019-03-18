Tea.context(function () {
	this.deleteReceiver = function (level, receiverId) {
		if (!window.confirm("确定要删除此接收人吗？")) {
			return;
		}
		this.$post("/agents/settings/deleteNoticeReceivers")
			.params({
				"agentId": this.agent.id,
				"level": level,
				"receiverId": receiverId
			})
			.refresh();
	};
});