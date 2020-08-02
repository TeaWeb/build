Tea.context(function () {
	this.deleteUser = function (userId) {
		if (!window.confirm("确定要删除此用户吗？")) {
			return;
		}
		this.$post(".acmeUserDelete")
			.params({
				"userId": userId
			})
			.refresh();
	};
});