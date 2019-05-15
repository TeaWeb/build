Tea.context(function () {
	this.connect = function () {
		this.error = "重试中...";
		this.$post(".connect").refresh();
	};

	/**
	 * 同步
	 */
	this.isSyncing = false;

	this.sync = function () {
		this.isSyncing = true;
		this.$post(".sync")
			.success(function () {
				this.$delay(function () {
					window.location.reload();
				}, 1000);
			});
	};
});