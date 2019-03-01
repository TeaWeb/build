Tea.context(function () {
	this.deleteGroup = function (groupId) {
		if (!window.confirm("确定要删除此分组吗？")) {
			return;
		}
		this.$post("/agents/groups/delete")
			.params({
				"groupId": groupId
			})
			.refresh();
	};
});