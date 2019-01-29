Tea.context(function () {
	this.deleteMedia = function (mediaId) {
		if (!window.confirm("确定要删除此媒介吗？")) {
			return;
		}
		this.$post("/notices/deleteMedia")
			.params({
				"mediaId": mediaId
			})
			.refresh();
	};
});