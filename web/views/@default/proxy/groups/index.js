Tea.context(function () {

	this.deleteGroup = function (groupId) {
		if (!window.confirm("确定要删除此分组吗？")) {
			return;
		}
		this.$post(".delete")
			.params({
				"serverId": this.server.id,
				"groupId": groupId
			})
			.refresh();
	};
});