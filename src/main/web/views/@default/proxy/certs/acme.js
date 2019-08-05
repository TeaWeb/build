Tea.context(function () {
	this.isRenewing = false;

	this.deleteTask = function (taskId) {
		if (!window.confirm("确定要删除此证书请求吗？正在使用的证书仍然会被保留")) {
			return;
		}
		this.$post(".acmeDeleteTask")
			.params({
				"taskId": taskId
			})
			.refresh();
	};

	this.renewTask = function (taskId) {
		if (!window.confirm("确定要现在更新吗？\n====================\n《严重注意》：Let's Encrypt对单个IP、单个域名请求频率有所限制，请勿频繁更新，建议在接近过期时更新。\n====================\n")) {
			return;
		}
		this.isRenewing = true;
		this.$post(".acmeRenewTask")
			.params({
				"taskId": taskId
			})
			.done(function () {
				this.isRenewing = false;
			})
			.refresh();
	};
});