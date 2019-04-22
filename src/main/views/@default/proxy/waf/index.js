Tea.context(function () {
	this.deleteWAF = function (wafId) {
		if (!window.confirm("确定要删除此WAF策略吗？")) {
			return;
		}
		this.$post(".delete")
			.params({
				"wafId": wafId
			})
			.refresh();
	};
});