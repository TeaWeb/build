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

	this.isRestoring = false;

	this.restore = function (file) {
		if (!window.confirm("确定要从此备份中恢复吗？")) {
			return
		}
		this.isRestoring = true;

		this.$post(".restore")
			.params({
				"file": file
			})
			.refresh();
	};

	this.clean = function () {
		if (!window.confirm("确定要清除30天以外的备份文件吗？")) {
			return;
		}
		this.$post(".clean")
			.success(function (resp) {
				if (resp.data.count > 0) {
					alert("成功清除了" + resp.data.count + "个备份文件");
				}
				window.location.reload();
			})
	};
});