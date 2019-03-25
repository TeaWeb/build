Tea.context(function () {
	/**
	 * 更多选项
	 */
	this.advancedOptionsVisible = false;

	this.showAdvancedOptions = function () {
		this.advancedOptionsVisible = !this.advancedOptionsVisible;
	};

	/**
	 * 启动
	 */
	this.startHttp = function () {
		if (!window.confirm("确定要启动此HTTP服务吗？")) {
			return;
		}
		this.$post(".startHttp")
			.params({
				"serverId": this.server.id,
			})
			.success(function () {
				window.location.reload();
			});
	};

	this.shutdownHttp = function () {
		if (!window.confirm("确定要关闭此HTTP服务吗？")) {
			return;
		}
		this.$post(".shutdownHttp")
			.params({
				"serverId": this.server.id,
			})
			.success(function () {
				window.location.reload();
			});
	};
});