Tea.context(function () {
	this.isCleaning = false;
	this.cleanResult = "";
	this.clean = function () {
		if (!window.confirm("确定要清理此策略关联的所有缓存吗？")) {
			return;
		}
		this.isCleaning = true;
		this.cleanResult = "";

		this.$post(".cleanPolicy")
			.params({
				"filename": this.policy.filename
			})
			.success(function (resp) {
				this.cleanResult = resp.data.result;
			})
			.fail(function () {
				this.cleanResult = resp.data.result;
			})
			.done(function () {
				this.isCleaning = false;
			})
	};
});