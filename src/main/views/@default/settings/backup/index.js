Tea.context(function () {
	this.backup = function () {
		if (!window.confirm("确定要立即备份配置吗？")) {
			return;
		}
		this.$post(".backup")
			.refresh();
	};

	this.deleteBackup = function (file) {
		if (!window.confirm("确定要删除此备份吗？")) {
			return;
		}
		this.$post(".delete")
			.params({
				"file": file
			})
			.refresh();
	};

	this.restore = function (file) {
		if (!window.confirm("确定要从此备份中恢复吗？")) {
			return
		}
		this.$post(".restore")
			.params({
				"file": file
			})
			.refresh();
	};
});