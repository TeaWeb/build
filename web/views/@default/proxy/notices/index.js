Tea.context(function () {
	this.levels.$each(function (k, v) {
		v.showGroupSetting = false;
	});

	this.deleteReceiver = function (level, receiverId) {
		if (!window.confirm("确定要删除此接收人吗？")) {
			return;
		}
		this.$post("/proxy/notices/deleteNoticeReceiver")
			.params({
				"serverId": this.server.id,
				"level": level,
				"receiverId": receiverId
			})
			.refresh();
	};

	this.showGroupSetting = function (level) {
		level.showGroupSetting = !level.showGroupSetting;
	};
});