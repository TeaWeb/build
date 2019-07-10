Tea.context(function () {
	this.deletePolicy = function (policyId) {
		if (!window.confirm("确定要删除此日志策略吗？")) {
			return;
		}

		this.$post(".delete")
			.params({
				"policyId": policyId
			})
			.refresh();
	};
});