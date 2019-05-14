Tea.context(function () {
	this.$delay(function () {
		this.$find("input[name='name']").focus();
	});

	this.submitSuccess = function () {
		if (!window.confirm("保存成功")) {
			return;
		}
		window.location = "/settings/cluster";
	};
});