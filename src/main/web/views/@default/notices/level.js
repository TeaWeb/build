Tea.context(function () {
	this.receivers.sort(function (r1, r2) {
		if (r1.mediaType > r2.mediaType) {
			return 1;
		}
		if (r1.mediaType < r2.mediaType) {
			return -1;
		}
		return 0;
	});

	this.deleteReceiver = function (receiverId) {
		if (!window.confirm("确定要删除此接收人吗？")) {
			return;
		}
		this.$post("/notices/deleteReceiver")
			.params({
				"level": this.level.code,
				"receiverId": receiverId
			})
			.refresh();
	};
});