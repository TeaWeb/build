Tea.context(function () {
	/**
	 * 更多选项
	 */
	this.advancedOptionsVisible = false;

	this.showAdvancedOptions = function () {
		this.advancedOptionsVisible = !this.advancedOptionsVisible;
	};

	/**
	 * 启动HTTP
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

	/**
	 * 关闭HTTP
	 */
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

	/**
	 * 启动TCP
	 */
	this.startTcp = function () {
		if (!window.confirm("确定要启动此TCP服务吗？")) {
			return;
		}
		this.$post(".startTcp")
			.params({
				"serverId": this.server.id,
			})
			.success(function () {
				window.location.reload();
			});
	};

	/**
	 * 关闭TCP
	 */
	this.shutdownTcp = function () {
		if (!window.confirm("确定要关闭此TCP服务吗？")) {
			return;
		}
		this.$post(".shutdownTcp")
			.params({
				"serverId": this.server.id,
			})
			.success(function () {
				window.location.reload();
			});
	};

	/**
	 * 通知设置
	 */
	this.noticeVisible = false;

	this.showNotice = function () {
		this.noticeVisible = !this.noticeVisible;
	};
});